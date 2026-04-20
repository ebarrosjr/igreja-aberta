using Jdb.Api.Controllers;
using Jdb.Api.Data;
using Jdb.Api.DTOs.Auth;
using Jdb.Api.Models;
using Jdb.Api.Services;
using Microsoft.AspNetCore.Mvc;
using Microsoft.EntityFrameworkCore;
using Xunit;

namespace Jdb.Api.Tests.Controllers
{
    public class AuthControllerTests : ControllerTestBase
    {
        private static AuthController CreateController(JdbContext context, IPasswordHasher passwordHasher, ITokenService tokenService)
        {
            return new AuthController(context, passwordHasher, tokenService, new TestWebHostEnvironment());
        }

        [Fact]
        public async Task Login_ReturnsJwtAndRefreshToken_WhenCredentialsAreValid()
        {
            await using JdbContext context = CreateContext();
            IPasswordHasher passwordHasher = new Pbkdf2PasswordHasher();
            ITokenService tokenService = new JwtTokenService(CreateConfiguration());
            Congregation congregation = await SeedCongregationAsync(context);
            await SeedUserAsync(context, passwordHasher, congregation.Id);
            AuthController controller = CreateController(context, passwordHasher, tokenService);

            ActionResult<AuthResponse> result = await controller.Login(new LoginRequest
            {
                Email = "user@example.com",
                Password = "Password123"
            });

            var ok = Assert.IsType<OkObjectResult>(result.Result);
            var response = Assert.IsType<AuthResponse>(ok.Value);
            Assert.False(string.IsNullOrWhiteSpace(response.Token));
            Assert.False(string.IsNullOrWhiteSpace(response.RefreshToken));
            Assert.Equal(1, await context.RefreshTokens.CountAsync());
        }

        [Fact]
        public async Task RefreshToken_ReturnsNewTokens_AndRevokesPreviousToken()
        {
            await using JdbContext context = CreateContext();
            IPasswordHasher passwordHasher = new Pbkdf2PasswordHasher();
            ITokenService tokenService = new JwtTokenService(CreateConfiguration());
            Congregation congregation = await SeedCongregationAsync(context);
            await SeedUserAsync(context, passwordHasher, congregation.Id);
            AuthController controller = CreateController(context, passwordHasher, tokenService);
            var login = Assert.IsType<OkObjectResult>((await controller.Login(new LoginRequest
            {
                Email = "user@example.com",
                Password = "Password123"
            })).Result);
            var loginResponse = Assert.IsType<AuthResponse>(login.Value);

            ActionResult<AuthResponse> result = await controller.RefreshToken(new RefreshTokenRequest
            {
                RefreshToken = loginResponse.RefreshToken
            });

            var ok = Assert.IsType<OkObjectResult>(result.Result);
            var response = Assert.IsType<AuthResponse>(ok.Value);
            Assert.NotEqual(loginResponse.RefreshToken, response.RefreshToken);
            Assert.Equal(2, await context.RefreshTokens.CountAsync());
            Assert.Equal(1, await context.RefreshTokens.CountAsync(t => t.RevokedAt != null));
        }

        [Fact]
        public async Task ForgotPassword_ReturnsAccepted_AndCreatesResetToken()
        {
            await using JdbContext context = CreateContext();
            IPasswordHasher passwordHasher = new Pbkdf2PasswordHasher();
            ITokenService tokenService = new JwtTokenService(CreateConfiguration());
            Congregation congregation = await SeedCongregationAsync(context);
            await SeedUserAsync(context, passwordHasher, congregation.Id);
            AuthController controller = CreateController(context, passwordHasher, tokenService);

            ActionResult<ForgotPasswordResponse> result = await controller.ForgotPassword(new ForgotPasswordRequest
            {
                Email = "user@example.com"
            });

            var accepted = Assert.IsType<AcceptedResult>(result.Result);
            var response = Assert.IsType<ForgotPasswordResponse>(accepted.Value);
            Assert.False(string.IsNullOrWhiteSpace(response.ResetToken));
            Assert.Equal(1, await context.PasswordResetTokens.CountAsync());
        }

        [Fact]
        public async Task ResetPassword_ChangesPassword_AndRevokesRefreshTokens()
        {
            await using JdbContext context = CreateContext();
            IPasswordHasher passwordHasher = new Pbkdf2PasswordHasher();
            ITokenService tokenService = new JwtTokenService(CreateConfiguration());
            Congregation congregation = await SeedCongregationAsync(context);
            User user = await SeedUserAsync(context, passwordHasher, congregation.Id);
            AuthController controller = CreateController(context, passwordHasher, tokenService);
            var resetToken = "reset-token";
            context.PasswordResetTokens.Add(new PasswordResetToken
            {
                UserId = user.Id,
                TokenHash = tokenService.HashToken(resetToken),
                ExpiresAt = DateTime.UtcNow.AddHours(1),
                CreatedAt = DateTime.UtcNow
            });
            context.RefreshTokens.Add(new RefreshToken
            {
                UserId = user.Id,
                TokenHash = tokenService.HashToken("refresh-token"),
                ExpiresAt = DateTime.UtcNow.AddDays(7),
                CreatedAt = DateTime.UtcNow
            });
            await context.SaveChangesAsync();

            IActionResult result = await controller.ResetPassword(new ResetPasswordRequest
            {
                Token = resetToken,
                NewPassword = "NewPassword123"
            });

            Assert.IsType<OkObjectResult>(result);
            User updatedUser = await context.Users.SingleAsync(u => u.Id == user.Id);
            Assert.True(passwordHasher.Verify("NewPassword123", updatedUser.Password));
            Assert.Equal(1, await context.RefreshTokens.CountAsync(t => t.RevokedAt != null));
        }
    }
}
