using Jdb.Api.Controllers;
using Jdb.Api.Data;
using Jdb.Api.DTOs;
using Jdb.Api.DTOs.Users;
using Jdb.Api.Models;
using Jdb.Api.Services;
using Microsoft.AspNetCore.Mvc;
using Microsoft.EntityFrameworkCore;
using Xunit;

namespace Jdb.Api.Tests.Controllers
{
    public class UsersControllerTests : ControllerTestBase
    {
        private static UsersController CreateController(JdbContext context, IPasswordHasher passwordHasher)
        {
            return new UsersController(context, passwordHasher);
        }

        [Fact]
        public async Task GetAll_ReturnsUsers()
        {
            await using JdbContext context = CreateContext();
            IPasswordHasher passwordHasher = new Pbkdf2PasswordHasher();
            Congregation congregation = await SeedCongregationAsync(context);
            await SeedUserAsync(context, passwordHasher, congregation.Id);
            UsersController controller = CreateController(context, passwordHasher);

            ActionResult<ApiResponse<IEnumerable<UserResponse>>> result = await controller.GetAll();

            var ok = Assert.IsType<ObjectResult>(result.Result);
            Assert.Equal(200, ok.StatusCode);
            var envelope = Assert.IsType<ApiResponse<IEnumerable<UserResponse>>>(ok.Value);
            var users = Assert.IsAssignableFrom<IEnumerable<UserResponse>>(envelope.Data);
            Assert.Single(users);
        }

        [Fact]
        public async Task Get_ReturnsUser_WhenUserExists()
        {
            await using JdbContext context = CreateContext();
            IPasswordHasher passwordHasher = new Pbkdf2PasswordHasher();
            Congregation congregation = await SeedCongregationAsync(context);
            User user = await SeedUserAsync(context, passwordHasher, congregation.Id);
            UsersController controller = CreateController(context, passwordHasher);

            ActionResult<ApiResponse<UserResponse>> result = await controller.Get(user.Id);

            var ok = Assert.IsType<ObjectResult>(result.Result);
            Assert.Equal(200, ok.StatusCode);
            var envelope = Assert.IsType<ApiResponse<UserResponse>>(ok.Value);
            UserResponse response = Assert.IsType<UserResponse>(envelope.Data);
            Assert.Equal(user.Email, response.Email);
        }

        [Fact]
        public async Task Create_AddsUser_WhenEmailIsAvailable()
        {
            await using JdbContext context = CreateContext();
            IPasswordHasher passwordHasher = new Pbkdf2PasswordHasher();
            Congregation congregation = await SeedCongregationAsync(context);
            UsersController controller = CreateController(context, passwordHasher);

            ActionResult<ApiResponse<UserResponse>> result = await controller.Create(new CreateUserRequest
            {
                CongregationId = congregation.Id,
                Name = "Novo Usuario",
                Email = "novo@example.com",
                Password = "Password123",
                Status = "active"
            });

            var created = Assert.IsType<ObjectResult>(result.Result);
            Assert.Equal(201, created.StatusCode);
            var envelope = Assert.IsType<ApiResponse<UserResponse>>(created.Value);
            UserResponse response = Assert.IsType<UserResponse>(envelope.Data);
            Assert.Equal("novo@example.com", response.Email);
            User storedUser = await context.Users.SingleAsync(u => u.Email == "novo@example.com");
            Assert.True(passwordHasher.Verify("Password123", storedUser.Password));
        }

        [Fact]
        public async Task Update_ChangesUser_AndKeepsPasswordWhenNotProvided()
        {
            await using JdbContext context = CreateContext();
            IPasswordHasher passwordHasher = new Pbkdf2PasswordHasher();
            Congregation congregation = await SeedCongregationAsync(context);
            User user = await SeedUserAsync(context, passwordHasher, congregation.Id);
            string originalPasswordHash = user.Password;
            UsersController controller = CreateController(context, passwordHasher);

            ActionResult<ApiResponse<UserResponse>> result = await controller.Update(user.Id, new UpdateUserRequest
            {
                CongregationId = congregation.Id,
                Name = "Usuario Atualizado",
                Email = "atualizado@example.com",
                Status = "active"
            });

            var ok = Assert.IsType<ObjectResult>(result.Result);
            Assert.Equal(200, ok.StatusCode);
            var envelope = Assert.IsType<ApiResponse<UserResponse>>(ok.Value);
            UserResponse response = Assert.IsType<UserResponse>(envelope.Data);
            Assert.Equal("Usuario Atualizado", response.Name);
            User storedUser = await context.Users.SingleAsync(u => u.Id == user.Id);
            Assert.Equal(originalPasswordHash, storedUser.Password);
        }

        [Fact]
        public async Task Delete_DeactivatesUser()
        {
            await using JdbContext context = CreateContext();
            IPasswordHasher passwordHasher = new Pbkdf2PasswordHasher();
            Congregation congregation = await SeedCongregationAsync(context);
            User user = await SeedUserAsync(context, passwordHasher, congregation.Id);
            UsersController controller = CreateController(context, passwordHasher);

            ActionResult<ApiResponse<object>> result = await controller.Delete(user.Id);

            var ok = Assert.IsType<ObjectResult>(result.Result);
            Assert.Equal(200, ok.StatusCode);
            var envelope = Assert.IsType<ApiResponse<object>>(ok.Value);
            Assert.Null(envelope.Data);
            User storedUser = await context.Users.SingleAsync(u => u.Id == user.Id);
            Assert.Equal("inactive", storedUser.Status);
        }
    }
}
