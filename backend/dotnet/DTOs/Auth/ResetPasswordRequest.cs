using System.ComponentModel.DataAnnotations;

namespace Jdb.Api.DTOs.Auth
{
    public class ResetPasswordRequest
    {
        [Required]
        public string Token { get; set; } = string.Empty;

        [Required]
        [MinLength(8)]
        [MaxLength(255)]
        public string NewPassword { get; set; } = string.Empty;
    }
}
