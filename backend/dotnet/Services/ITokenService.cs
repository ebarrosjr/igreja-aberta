using Jdb.Api.DTOs.Auth;
using Jdb.Api.Models;

namespace Jdb.Api.Services
{
    public interface ITokenService
    {
        AuthResponse CreateAuthResponse(User user);
        string GenerateOpaqueToken();
        string HashToken(string token);
    }
}
