using System.Text;
using Jdb.Api.Data;
using Jdb.Api.DTOs;
using Jdb.Api.Middleware;
using Jdb.Api.Seeders;
using Jdb.Api.Services;
using Microsoft.AspNetCore.Authentication.JwtBearer;
using Microsoft.AspNetCore.Mvc;
using Microsoft.EntityFrameworkCore;
using Microsoft.IdentityModel.Tokens;

var builder = WebApplication.CreateBuilder(args);

builder.Services.AddDbContext<JdbContext>(options =>
    options.UseMySql(
        builder.Configuration.GetConnectionString("DefaultConnection"),
        new MySqlServerVersion(new Version(8, 0, 36))));

builder.Services.AddScoped<IPasswordHasher, BcryptPasswordHasher>();
builder.Services.AddScoped<ITokenService, JwtTokenService>();
builder.Services.AddControllers()
    .ConfigureApiBehaviorOptions(options =>
    {
        options.InvalidModelStateResponseFactory = context =>
        {
            string message = string.Join(" ", context.ModelState.Values
                .SelectMany(v => v.Errors)
                .Select(e => string.IsNullOrWhiteSpace(e.ErrorMessage) ? "Requisicao invalida." : e.ErrorMessage));

            return new BadRequestObjectResult(new ApiResponse<object>
            {
                Code = StatusCodes.Status400BadRequest,
                Data = null,
                Message = string.IsNullOrWhiteSpace(message) ? "Requisicao invalida." : message
            });
        };
    });
builder.Services.AddEndpointsApiExplorer();
builder.Services.AddSwaggerGen();
builder.Services.AddOpenApi();
builder.Services
    .AddAuthentication(JwtBearerDefaults.AuthenticationScheme)
    .AddJwtBearer(options =>
    {
        string jwtKey = builder.Configuration["Jwt:Key"]
            ?? throw new InvalidOperationException("Jwt:Key nao configurado.");

        options.TokenValidationParameters = new TokenValidationParameters
        {
            ValidateIssuer = true,
            ValidateAudience = true,
            ValidateLifetime = true,
            ValidateIssuerSigningKey = true,
            ValidIssuer = builder.Configuration["Jwt:Issuer"],
            ValidAudience = builder.Configuration["Jwt:Audience"],
            IssuerSigningKey = new SymmetricSecurityKey(Encoding.UTF8.GetBytes(jwtKey)),
            ClockSkew = TimeSpan.FromMinutes(1)
        };

        options.Events = new JwtBearerEvents
        {
            OnChallenge = async context =>
            {
                context.HandleResponse();
                context.Response.StatusCode = StatusCodes.Status401Unauthorized;
                context.Response.ContentType = "application/json";
                await context.Response.WriteAsJsonAsync(new ApiResponse<object>
                {
                    Code = StatusCodes.Status401Unauthorized,
                    Data = null,
                    Message = "Token ausente ou invalido."
                });
            },
            OnForbidden = async context =>
            {
                context.Response.StatusCode = StatusCodes.Status403Forbidden;
                context.Response.ContentType = "application/json";
                await context.Response.WriteAsJsonAsync(new ApiResponse<object>
                {
                    Code = StatusCodes.Status403Forbidden,
                    Data = null,
                    Message = "Acesso negado."
                });
            }
        };
    });
builder.Services.AddAuthorization();

var app = builder.Build();

bool swaggerEnabled = app.Environment.IsDevelopment()
    || builder.Configuration.GetValue<bool>("Swagger:Enabled");

if (swaggerEnabled)
{
    app.MapOpenApi();
    app.UseSwagger();
    app.UseSwaggerUI();
}

app.UseMiddleware<ExceptionHandlingMiddleware>();
app.UseHttpsRedirection();
app.UseStatusCodePages(async context =>
{
    HttpResponse response = context.HttpContext.Response;
    if (response.HasStarted || response.ContentLength.HasValue || response.StatusCode < 400)
    {
        return;
    }

    response.ContentType = "application/json";
    await response.WriteAsJsonAsync(new ApiResponse<object>
    {
        Code = response.StatusCode,
        Data = null,
        Message = response.StatusCode == StatusCodes.Status404NotFound
            ? "Recurso nao encontrado."
            : "Falha ao processar a requisicao."
    });
});
app.UseAuthentication();
app.UseAuthorization();
app.MapControllers();

await app.SeedInitialDataAsync();

app.Run();
