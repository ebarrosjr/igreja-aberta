using Jdb.Api.Data;
using Jdb.Api.Models;
using Jdb.Api.Services;
using Microsoft.EntityFrameworkCore;

namespace Jdb.Api.Seeders
{
    public static class InitialDataSeeder
    {
        public static async Task SeedInitialDataAsync(this WebApplication app)
        {
            using IServiceScope scope = app.Services.CreateScope();
            var context = scope.ServiceProvider.GetRequiredService<JdbContext>();
            var passwordHasher = scope.ServiceProvider.GetRequiredService<IPasswordHasher>();
            var configuration = scope.ServiceProvider.GetRequiredService<IConfiguration>();

            Congregation congregation = await context.Congregations
                .FirstOrDefaultAsync(c => c.Email == "admin@igrejaaberta.local")
                ?? await CreateCongregationAsync(context);

            bool adminExists = await context.Users.AnyAsync(u => u.Email == "ebarrosjr@gmail.com");
            if (!adminExists)
            {
                DateTime now = DateTime.UtcNow;
                string adminPassword = configuration["Seed:AdminPassword"] ?? "Admin@123456";

                context.Users.Add(new User
                {
                    CongregationId = congregation.Id,
                    Name = "Administrador",
                    Email = "ebarrosjr@gmail.com",
                    Password = passwordHasher.Hash(adminPassword),
                    Status = "active",
                    CreatedAt = now,
                    UpdatedAt = now
                });

                await context.SaveChangesAsync();
            }
        }

        private static async Task<Congregation> CreateCongregationAsync(JdbContext context)
        {
            DateTime now = DateTime.UtcNow;
            var congregation = new Congregation
            {
                Name = "Congregacao Inicial",
                Description = "Congregacao criada automaticamente pelo seeder inicial.",
                Email = "admin@igrejaaberta.local",
                Status = "active",
                CreatedAt = now,
                UpdatedAt = now
            };

            context.Congregations.Add(congregation);
            await context.SaveChangesAsync();

            return congregation;
        }
    }
}
