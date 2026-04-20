using Jdb.Api.DTOs;
using Microsoft.AspNetCore.Mvc;

namespace Jdb.Api.Controllers
{
    public abstract class ApiControllerBase : ControllerBase
    {
        protected ActionResult<ApiResponse<T>> ApiOk<T>(T? data, string message = "OK")
        {
            return StatusCode(StatusCodes.Status200OK, BuildResponse(StatusCodes.Status200OK, data, message));
        }

        protected ActionResult<ApiResponse<T>> ApiCreated<T>(T? data, string message = "Criado com sucesso")
        {
            return StatusCode(StatusCodes.Status201Created, BuildResponse(StatusCodes.Status201Created, data, message));
        }

        protected ActionResult<ApiResponse<T>> ApiAccepted<T>(T? data, string message = "Aceito")
        {
            return StatusCode(StatusCodes.Status202Accepted, BuildResponse(StatusCodes.Status202Accepted, data, message));
        }

        protected ActionResult<ApiResponse<T>> ApiBadRequest<T>(string message)
        {
            return StatusCode(StatusCodes.Status400BadRequest, BuildResponse<T>(StatusCodes.Status400BadRequest, default, message));
        }

        protected ActionResult<ApiResponse<T>> ApiUnauthorized<T>(string message)
        {
            return StatusCode(StatusCodes.Status401Unauthorized, BuildResponse<T>(StatusCodes.Status401Unauthorized, default, message));
        }

        protected ActionResult<ApiResponse<T>> ApiNotFound<T>(string message)
        {
            return StatusCode(StatusCodes.Status404NotFound, BuildResponse<T>(StatusCodes.Status404NotFound, default, message));
        }

        protected ActionResult<ApiResponse<T>> ApiConflict<T>(string message)
        {
            return StatusCode(StatusCodes.Status409Conflict, BuildResponse<T>(StatusCodes.Status409Conflict, default, message));
        }

        protected static ApiResponse<T> BuildResponse<T>(int code, T? data, string message)
        {
            return new ApiResponse<T>
            {
                Code = code,
                Data = data,
                Message = message
            };
        }
    }
}
