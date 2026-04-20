using System.ComponentModel.DataAnnotations;

namespace Jdb.Api.DTOs.Auth
{
    public class RefreshTokenRequest
    {
        [Required]
        public string RefreshToken { get; set; } = string.Empty;
    }
}
