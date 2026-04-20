using System.ComponentModel.DataAnnotations;

namespace Jdb.Api.Models
{
    public class RefreshToken
    {
        [Key]
        public long Id { get; set; }
        public long UserId { get; set; }
        [Required]
        [MaxLength(128)]
        public string TokenHash { get; set; } = string.Empty;
        public DateTime ExpiresAt { get; set; }
        public DateTime? RevokedAt { get; set; }
        public DateTime CreatedAt { get; set; }
        public User? User { get; set; }
    }
}
