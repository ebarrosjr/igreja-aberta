using Jdb.Api.Data;
using Jdb.Api.Models;
using Jdb.Api.Services;
using Microsoft.AspNetCore.Hosting;
using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.FileProviders;
using Microsoft.Extensions.Hosting;

namespace Jdb.Api.Tests
{
    public abstract class ControllerTestBase
    {
        protected static JdbContext CreateContext()
        {
            var options = new DbContextOptionsBuilder<JdbContext>()
                .UseInMemoryDatabase(Guid.NewGuid().ToString())
                .Options;

            return new JdbContext(options);
        }

        protected static IConfiguration CreateConfiguration()
        {
            return new ConfigurationBuilder()
                .AddInMemoryCollection(new Dictionary<string, string?>
                {
                    ["Jwt:Issuer"] = "igreja-aberta-tests",
                    ["Jwt:Audience"] = "igreja-aberta-tests-api",
                    ["Jwt:Key"] = "TEST_SECRET_KEY_WITH_MORE_THAN_32_BYTES",
                    ["Jwt:AccessTokenMinutes"] = "15",
                    ["Jwt:RefreshTokenDays"] = "7"
                })
                .Build();
        }

        protected static async Task<Congregation> SeedCongregationAsync(JdbContext context)
        {
            var congregation = new Congregation
            {
                Name = "Congregacao Teste",
                Email = "teste@igreja.local",
                Status = "active",
                CreatedAt = DateTime.UtcNow,
                UpdatedAt = DateTime.UtcNow
            };

            context.Congregations.Add(congregation);
            await context.SaveChangesAsync();
            return congregation;
        }

        protected static async Task<User> SeedUserAsync(
            JdbContext context,
            IPasswordHasher passwordHasher,
            long congregationId,
            string email = "user@example.com",
            string password = "Password123")
        {
            var user = new User
            {
                CongregationId = congregationId,
                Name = "Usuario Teste",
                Email = email,
                Password = passwordHasher.Hash(password),
                Status = "active",
                CreatedAt = DateTime.UtcNow,
                UpdatedAt = DateTime.UtcNow
            };

            context.Users.Add(user);
            await context.SaveChangesAsync();
            return user;
        }
    }

    public sealed class TestWebHostEnvironment : IWebHostEnvironment
    {
        public string ApplicationName { get; set; } = "Jdb.Api.Tests";
        public IFileProvider ContentRootFileProvider { get; set; } = new NullFileProvider();
        public string ContentRootPath { get; set; } = Directory.GetCurrentDirectory();
        public string EnvironmentName { get; set; } = Environments.Development;
        public string WebRootPath { get; set; } = Directory.GetCurrentDirectory();
        public IFileProvider WebRootFileProvider { get; set; } = new NullFileProvider();
    }
}
