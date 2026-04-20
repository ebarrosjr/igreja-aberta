namespace Jdb.Api.DTOs
{
    public class ApiResponse<T>
    {
        public int Code { get; set; }
        public T? Data { get; set; }
        public string Message { get; set; } = string.Empty;
        public string Stack { get; set; } = "dotnet";
        public string Db { get; set; } = "mysql";
    }
}
