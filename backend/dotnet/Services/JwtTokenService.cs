using System.IdentityModel.Tokens.Jwt;
using System.Security.Claims;
using System.Security.Cryptography;
using System.Text;
using Jdb.Api.DTOs.Auth;
using Jdb.Api.Models;
using Microsoft.IdentityModel.Tokens;

namespace Jdb.Api.Services
{
    public class JwtTokenService : ITokenService
    {
        private readonly IConfiguration _configuration;

        public JwtTokenService(IConfiguration configuration)
        {
            _configuration = configuration;
        }

        public AuthResponse CreateAuthResponse(User user)
        {
            DateTime expiresAt = DateTime.UtcNow.AddMinutes(GetAccessTokenMinutes());
            string refreshToken = GenerateOpaqueToken();

            var claims = new List<Claim>
            {
                new(JwtRegisteredClaimNames.Sub, user.Id.ToString()),
                new(JwtRegisteredClaimNames.Email, user.Email),
                new(ClaimTypes.NameIdentifier, user.Id.ToString()),
                new(ClaimTypes.Email, user.Email),
                new("congregation_id", user.CongregationId.ToString())
            };

            var credentials = new SigningCredentials(
                new SymmetricSecurityKey(Encoding.UTF8.GetBytes(GetJwtKey())),
                SecurityAlgorithms.HmacSha256);

            var jwt = new JwtSecurityToken(
                issuer: _configuration["Jwt:Issuer"],
                audience: _configuration["Jwt:Audience"],
                claims: claims,
                expires: expiresAt,
                signingCredentials: credentials);

            return new AuthResponse
            {
                Token = new JwtSecurityTokenHandler().WriteToken(jwt),
                ExpiresAt = expiresAt,
                RefreshToken = refreshToken,
                RefreshTokenExpiresAt = DateTime.UtcNow.AddDays(GetRefreshTokenDays())
            };
        }

        public string GenerateOpaqueToken()
        {
            return Convert.ToBase64String(RandomNumberGenerator.GetBytes(64));
        }

        public string HashToken(string token)
        {
            byte[] hash = SHA256.HashData(Encoding.UTF8.GetBytes(token));
            return Convert.ToHexString(hash).ToLowerInvariant();
        }

        private string GetJwtKey()
        {
            string? key = _configuration["Jwt:Key"];
            if (string.IsNullOrWhiteSpace(key) || Encoding.UTF8.GetByteCount(key) < 32)
            {
                throw new InvalidOperationException("Jwt:Key deve ter pelo menos 32 bytes.");
            }

            return key;
        }

        private int GetAccessTokenMinutes()
        {
            return int.TryParse(_configuration["Jwt:AccessTokenMinutes"], out int minutes) ? minutes : 15;
        }

        private int GetRefreshTokenDays()
        {
            return int.TryParse(_configuration["Jwt:RefreshTokenDays"], out int days) ? days : 7;
        }
    }
}
