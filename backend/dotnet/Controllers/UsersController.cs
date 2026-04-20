using Jdb.Api.Data;
using Jdb.Api.DTOs.Users;
using Jdb.Api.Models;
using Jdb.Api.Services;
using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Mvc;
using Microsoft.EntityFrameworkCore;

namespace Jdb.Api.Controllers
{
    [ApiController]
    [Authorize]
    [Route("api/[controller]")]
    public class UsersController : ControllerBase
    {
        private readonly JdbContext _context;
        private readonly IPasswordHasher _passwordHasher;

        public UsersController(JdbContext context, IPasswordHasher passwordHasher)
        {
            _context = context;
            _passwordHasher = passwordHasher;
        }

        [HttpGet]
        public async Task<ActionResult<IEnumerable<UserResponse>>> GetAll()
        {
            List<UserResponse> users = await _context.Users
                .AsNoTracking()
                .OrderBy(u => u.Name)
                .Select(u => new UserResponse
                {
                    Id = u.Id,
                    CongregationId = u.CongregationId,
                    Name = u.Name,
                    Email = u.Email,
                    Status = u.Status,
                    LastLoginAt = u.LastLoginAt,
                    CreatedAt = u.CreatedAt,
                    UpdatedAt = u.UpdatedAt
                })
                .ToListAsync();

            return Ok(users);
        }

        [HttpGet("{id:long}")]
        public async Task<ActionResult<UserResponse>> Get(long id)
        {
            User? user = await _context.Users.AsNoTracking().FirstOrDefaultAsync(u => u.Id == id);
            if (user is null)
            {
                return NotFound(new { message = "Usuario nao encontrado." });
            }

            return Ok(ToResponse(user));
        }

        [HttpPost]
        public async Task<ActionResult<UserResponse>> Create(CreateUserRequest request)
        {
            string email = request.Email.Trim().ToLowerInvariant();

            if (!await _context.Congregations.AnyAsync(c => c.Id == request.CongregationId))
            {
                return BadRequest(new { message = "Congregacao informada nao existe." });
            }

            if (await _context.Users.AnyAsync(u => u.Email == email))
            {
                return Conflict(new { message = "E-mail ja cadastrado." });
            }

            DateTime now = DateTime.UtcNow;
            var user = new User
            {
                CongregationId = request.CongregationId,
                Name = request.Name.Trim(),
                Email = email,
                Password = _passwordHasher.Hash(request.Password),
                Status = request.Status.Trim().ToLowerInvariant(),
                CreatedAt = now,
                UpdatedAt = now
            };

            _context.Users.Add(user);
            await _context.SaveChangesAsync();

            return CreatedAtAction(nameof(Get), new { id = user.Id }, ToResponse(user));
        }

        [HttpPut("{id:long}")]
        public async Task<ActionResult<UserResponse>> Update(long id, UpdateUserRequest request)
        {
            User? user = await _context.Users.FirstOrDefaultAsync(u => u.Id == id);
            if (user is null)
            {
                return NotFound(new { message = "Usuario nao encontrado." });
            }

            string email = request.Email.Trim().ToLowerInvariant();
            if (!await _context.Congregations.AnyAsync(c => c.Id == request.CongregationId))
            {
                return BadRequest(new { message = "Congregacao informada nao existe." });
            }

            if (await _context.Users.AnyAsync(u => u.Id != id && u.Email == email))
            {
                return Conflict(new { message = "E-mail ja cadastrado." });
            }

            user.CongregationId = request.CongregationId;
            user.Name = request.Name.Trim();
            user.Email = email;
            user.Status = request.Status.Trim().ToLowerInvariant();
            user.UpdatedAt = DateTime.UtcNow;

            if (!string.IsNullOrWhiteSpace(request.Password))
            {
                user.Password = _passwordHasher.Hash(request.Password);
            }

            await _context.SaveChangesAsync();

            return Ok(ToResponse(user));
        }

        [HttpDelete("{id:long}")]
        public async Task<IActionResult> Delete(long id)
        {
            User? user = await _context.Users.FirstOrDefaultAsync(u => u.Id == id);
            if (user is null)
            {
                return NotFound(new { message = "Usuario nao encontrado." });
            }

            user.Status = "inactive";
            user.UpdatedAt = DateTime.UtcNow;
            await _context.SaveChangesAsync();

            return NoContent();
        }

        private static UserResponse ToResponse(User user)
        {
            return new UserResponse
            {
                Id = user.Id,
                CongregationId = user.CongregationId,
                Name = user.Name,
                Email = user.Email,
                Status = user.Status,
                LastLoginAt = user.LastLoginAt,
                CreatedAt = user.CreatedAt,
                UpdatedAt = user.UpdatedAt
            };
        }
    }
}
