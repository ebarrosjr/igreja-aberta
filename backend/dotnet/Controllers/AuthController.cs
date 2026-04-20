using Jdb.Api.Data;
using Jdb.Api.DTOs.Auth;
using Jdb.Api.Models;
using Jdb.Api.Services;
using Microsoft.AspNetCore.Mvc;
using Microsoft.EntityFrameworkCore;

namespace Jdb.Api.Controllers
{
    [ApiController]
    [Route("api/[controller]")]
    public class AuthController : ControllerBase
    {
        private readonly JdbContext _context;
        private readonly IPasswordHasher _passwordHasher;
        private readonly ITokenService _tokenService;
        private readonly IWebHostEnvironment _environment;

        public AuthController(
            JdbContext context,
            IPasswordHasher passwordHasher,
            ITokenService tokenService,
            IWebHostEnvironment environment)
        {
            _context = context;
            _passwordHasher = passwordHasher;
            _tokenService = tokenService;
            _environment = environment;
        }

        [HttpPost("login")]
        public async Task<ActionResult<AuthResponse>> Login(LoginRequest request)
        {
            string email = request.Email.Trim().ToLowerInvariant();
            User? user = await _context.Users.FirstOrDefaultAsync(u => u.Email == email);

            if (user is null || user.Status != "active" || !_passwordHasher.Verify(request.Password, user.Password))
            {
                return Unauthorized(new { message = "Credenciais invalidas." });
            }

            AuthResponse response = _tokenService.CreateAuthResponse(user);
            DateTime now = DateTime.UtcNow;

            user.LastLoginAt = now;
            user.UpdatedAt = now;
            _context.RefreshTokens.Add(new RefreshToken
            {
                UserId = user.Id,
                TokenHash = _tokenService.HashToken(response.RefreshToken),
                ExpiresAt = response.RefreshTokenExpiresAt,
                CreatedAt = now
            });

            await _context.SaveChangesAsync();

            return Ok(response);
        }

        [HttpPost("refresh-token")]
        public async Task<ActionResult<AuthResponse>> RefreshToken(RefreshTokenRequest request)
        {
            string tokenHash = _tokenService.HashToken(request.RefreshToken);
            RefreshToken? currentRefreshToken = await _context.RefreshTokens
                .Include(t => t.User)
                .FirstOrDefaultAsync(t => t.TokenHash == tokenHash);

            if (currentRefreshToken is null
                || currentRefreshToken.User is null
                || currentRefreshToken.User.Status != "active"
                || currentRefreshToken.RevokedAt is not null
                || currentRefreshToken.ExpiresAt <= DateTime.UtcNow)
            {
                return Unauthorized(new { message = "Refresh token invalido." });
            }

            AuthResponse response = _tokenService.CreateAuthResponse(currentRefreshToken.User);
            DateTime now = DateTime.UtcNow;

            currentRefreshToken.RevokedAt = now;
            _context.RefreshTokens.Add(new RefreshToken
            {
                UserId = currentRefreshToken.UserId,
                TokenHash = _tokenService.HashToken(response.RefreshToken),
                ExpiresAt = response.RefreshTokenExpiresAt,
                CreatedAt = now
            });

            await _context.SaveChangesAsync();

            return Ok(response);
        }

        [HttpPost("forgot-password")]
        public async Task<ActionResult<ForgotPasswordResponse>> ForgotPassword(ForgotPasswordRequest request)
        {
            string email = request.Email.Trim().ToLowerInvariant();
            User? user = await _context.Users.FirstOrDefaultAsync(u => u.Email == email && u.Status == "active");

            if (user is null)
            {
                return Accepted(new ForgotPasswordResponse());
            }

            string resetToken = _tokenService.GenerateOpaqueToken();
            DateTime expiresAt = DateTime.UtcNow.AddHours(1);

            _context.PasswordResetTokens.Add(new PasswordResetToken
            {
                UserId = user.Id,
                TokenHash = _tokenService.HashToken(resetToken),
                ExpiresAt = expiresAt,
                CreatedAt = DateTime.UtcNow
            });

            await _context.SaveChangesAsync();

            return Accepted(new ForgotPasswordResponse
            {
                ResetToken = _environment.IsDevelopment() ? resetToken : null,
                ExpiresAt = _environment.IsDevelopment() ? expiresAt : null
            });
        }

        [HttpPost("reset-password")]
        public async Task<IActionResult> ResetPassword(ResetPasswordRequest request)
        {
            string tokenHash = _tokenService.HashToken(request.Token);
            PasswordResetToken? resetToken = await _context.PasswordResetTokens
                .Include(t => t.User)
                .FirstOrDefaultAsync(t => t.TokenHash == tokenHash);

            if (resetToken is null
                || resetToken.User is null
                || resetToken.UsedAt is not null
                || resetToken.ExpiresAt <= DateTime.UtcNow)
            {
                return BadRequest(new { message = "Token de redefinicao invalido ou expirado." });
            }

            DateTime now = DateTime.UtcNow;
            resetToken.UsedAt = now;
            resetToken.User.Password = _passwordHasher.Hash(request.NewPassword);
            resetToken.User.UpdatedAt = now;

            List<RefreshToken> activeRefreshTokens = await _context.RefreshTokens
                .Where(t => t.UserId == resetToken.UserId && t.RevokedAt == null)
                .ToListAsync();

            foreach (RefreshToken refreshToken in activeRefreshTokens)
            {
                refreshToken.RevokedAt = now;
            }

            await _context.SaveChangesAsync();

            return Ok(new { message = "Senha redefinida com sucesso." });
        }
    }
}
