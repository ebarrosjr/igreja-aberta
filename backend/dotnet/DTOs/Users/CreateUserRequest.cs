using System.ComponentModel.DataAnnotations;

namespace Jdb.Api.DTOs.Users
{
    public class CreateUserRequest
    {
        [Required]
        public long CongregationId { get; set; }

        [Required]
        [MaxLength(255)]
        public string Name { get; set; } = string.Empty;

        [Required]
        [EmailAddress]
        [MaxLength(255)]
        public string Email { get; set; } = string.Empty;

        [Required]
        [MinLength(8)]
        [MaxLength(255)]
        public string Password { get; set; } = string.Empty;

        [Required]
        [MaxLength(50)]
        public string Status { get; set; } = "active";
    }
}
