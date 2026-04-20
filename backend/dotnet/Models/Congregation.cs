using System.ComponentModel.DataAnnotations;

namespace Jdb.Api.Models
{
    public class Congregation
    {
        [Key]
        public long Id { get; set; }
        [Required]
        [MaxLength(255)]
        public string Name { get; set; } = string.Empty;
        [MaxLength(1000)]
        public string? Description { get; set; }
        [MaxLength(30)]
        public string? Phone { get; set; }
        [MaxLength(255)]
        public string? Email { get; set; }
        [MaxLength(20)]
        public string? ZipCode { get; set; }
        [MaxLength(255)]
        public string? Address { get; set; }
        [MaxLength(20)]
        public string? Number { get; set; }
        [MaxLength(255)]
        public string? Complement { get; set; }
        [MaxLength(255)]
        public string? Neighborhood { get; set; }
        [MaxLength(255)]
        public string? City { get; set; }
        [MaxLength(100)]
        public string? State { get; set; }
        [Required]
        [MaxLength(50)]
        public string Status { get; set; } = "active";
        public DateTime CreatedAt { get; set; }
        public DateTime UpdatedAt { get; set; }
        public ICollection<User> Users { get; set; } = new List<User>();
    }

}
