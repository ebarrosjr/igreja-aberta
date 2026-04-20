using Jdb.Api.Data;
using Jdb.Api.DTOs;
using Jdb.Api.DTOs.Auth;
using Jdb.Api.Models;
using Jdb.Api.Services;
using Microsoft.AspNetCore.Mvc;
using Microsoft.EntityFrameworkCore;

namespace Jdb.Api.Controllers
{
    [ApiController]
    [Route("api/[controller]")]
    public class AuthController : ApiControllerBase
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
        public async Task<ActionResult<ApiResponse<AuthResponse>>> Login(LoginRequest request)
        {
            string email = request.Email.Trim().ToLowerInvariant();
            User? user = await _context.Users.FirstOrDefaultAsync(u => u.Email == email);

            if (user is null || user.Status != "active" || !_passwordHasher.Verify(request.Password, user.Password))
            {
                return ApiUnauthorized<AuthResponse>("Credenciais invalidas.");
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

            return ApiOk(response, "Login realizado com sucesso.");
        }

        [HttpPost("refresh-token")]
        public async Task<ActionResult<ApiResponse<AuthResponse>>> RefreshToken(RefreshTokenRequest request)
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
                return ApiUnauthorized<AuthResponse>("Refresh token invalido.");
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

            return ApiOk(response, "Token atualizado com sucesso.");
        }

        [HttpPost("forgot-password")]
        public async Task<ActionResult<ApiResponse<ForgotPasswordResponse>>> ForgotPassword(ForgotPasswordRequest request)
        {
            string email = request.Email.Trim().ToLowerInvariant();
            User? user = await _context.Users.FirstOrDefaultAsync(u => u.Email == email && u.Status == "active");

            if (user is null)
            {
                return ApiAccepted(new ForgotPasswordResponse(), "Se o e-mail existir, as instrucoes de redefinicao serao enviadas.");
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

            return ApiAccepted(new ForgotPasswordResponse
            {
                ResetToken = _environment.IsDevelopment() ? resetToken : null,
                ExpiresAt = _environment.IsDevelopment() ? expiresAt : null
            }, "Se o e-mail existir, as instrucoes de redefinicao serao enviadas.");
        }

        [HttpPost("reset-password")]
        public async Task<ActionResult<ApiResponse<object>>> ResetPassword(ResetPasswordRequest request)
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
                return ApiBadRequest<object>("Token de redefinicao invalido ou expirado.");
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

            return ApiOk<object>(null, "Senha redefinida com sucesso.");
        }
    }
}
