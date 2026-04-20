namespace Jdb.Api.DTOs.Auth
{
    public class ForgotPasswordResponse
    {
        public string Message { get; set; } = "Se o e-mail existir, as instrucoes de redefinicao serao enviadas.";
        public string? ResetToken { get; set; }
        public DateTime? ExpiresAt { get; set; }
    }
}
