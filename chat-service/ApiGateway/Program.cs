
using Microsoft.AspNetCore.Builder;     // Import ASP.NET Core middleware and application builder
using Microsoft.Extensions.DependencyInjection; // Import dependency injection services
using System.Net.Http;
using Serilog;
var builder = WebApplication.CreateBuilder(args);

// Add services to the container.
// Learn more about configuring Swagger/OpenAPI at https://aka.ms/aspnetcore/swashbuckle
builder.Services.AddEndpointsApiExplorer();
builder.Services.AddSwaggerGen();

builder.Services.AddHttpClient("auth-service", client => // Register an HTTP client for auth-service
{                                                       // This is like setting up a client in Go with `http.Client`.
client.BaseAddress = new Uri("http://localhost:2021");    // In Go, you'd set this in a `http.Client` with a custom transport or base URL.
    client.DefaultRequestHeaders.Add("Accept", "application/json"); // Add a default header (optional)
});

var app = builder.Build();
app.UseRouting();
// Configure the HTTP request pipeline.
// app.Map("auth/register", async (HttpContext context, IHttpClientFactory clientFactory, ILogger<Program>log) =>
// {
//     var client = clientFactory.CreateClient("auth-service");
//     var path = context.Request.Path.Value;
//     var query = context.Request.QueryString.Value;
//     log.LogInformation("proxying:Path={path}", path);
//    if (path == "/auth/register")
//     {
//         var backendEndpoint = path.Replace("/auth", "");
//         if (backendEndpoint == "/register")
//         {
//              string cleanBackendpath = backendEndpoint;
//             if (backendEndpoint.Length > 0 && backendEndpoint[0] == '/')
//             {
//                 cleanBackendpath = backendEndpoint.Length > 1 ? backendEndpoint.Substring(1) : " ";
//             }
//             {

//  }
//             var request = new HttpRequestMessage
//             {
//                 Method = new HttpMethod(context.Request.Method),
//                 RequestUri = new Uri(client.BaseAddress + cleanBackendpath + query) // Construct the full URL
//                                                                                    //   http://auth-service:8080/register
//                                                                                                   //   http://localhost:2021/auth/register/



//             }; log.LogInformation("Request URI={RequestURI}", request.RequestUri);
//             request.Content = new StreamContent(context.Request.Body);
//            using var response = await client.SendAsync(request, HttpCompletionOption.ResponseHeadersRead, context.RequestAborted);
//         // Send the request to auth-service and await the response
//         // In Go, this is like `resp, err := http.DefaultClient.Do(req)`, but async with `await`.
//         context.Response.StatusCode = (int)response.StatusCode; // Set the response status code
//                                                                 // Similar to setting `w.WriteHeader` in Go's `http.ResponseWriter`.
//         await response.Content.CopyToAsync(context.Response.Body); // Copy the response body to the client
//                                                                    // In Go, you'd use `io.Copy(w, resp.Body)` to write the response.
//          }
//         ;



//     }
//     ;
// });
app.Map("auth/register", async (
        HttpContext context,
        IHttpClientFactory clientFactory,
        ILogger<Program> log) =>
{
    var client = clientFactory.CreateClient("auth-service");
    var path = context.Request.Path.Value;
    var query = context.Request.QueryString.Value;

    // log.LogInformation("proxying:Path={path}", path);

    if (path == "/auth/register")
    {
        var backendEndpoint = path.Replace("/auth", "");   // -> "/register"
        var cleanBackendPath = backendEndpoint.StartsWith("/")
                               ? backendEndpoint[1..]      // remove leading /
                               : backendEndpoint;

        var request = new HttpRequestMessage
        {
            Method = new HttpMethod(context.Request.Method),
            RequestUri = new Uri(client.BaseAddress + cleanBackendPath + query)
        };

        // *** copy the body ***
        request.Content = new StreamContent(context.Request.Body);

        // *** copy headers ***
        foreach (var h in context.Request.Headers)
            if (!request.Headers.TryAddWithoutValidation(h.Key, h.Value.ToArray()))
                request.Content?.Headers.TryAddWithoutValidation(h.Key, h.Value.ToArray());

        // log.LogInformation("Request URI={RequestURI}", request.RequestUri);

        using var response = await client.SendAsync(
                                request,
                                HttpCompletionOption.ResponseHeadersRead,
                                context.RequestAborted);

        context.Response.StatusCode = (int)response.StatusCode;
        await response.Content.CopyToAsync(context.Response.Body);
    }
});
app.Map("auth/login", async (
        HttpContext context,
        IHttpClientFactory clientFactory,
        ILogger<Program> log) =>
{
    var client = clientFactory.CreateClient("auth-service");
    var path = context.Request.Path.Value;
    var query = context.Request.QueryString.Value;

    // log.LogInformation("proxying:Path={path}", path);

    if (path == "/auth/login")
    {
        var backendEndpoint = path.Replace("/auth", "");   // -> "/login"
        var cleanBackendPath = backendEndpoint.StartsWith("/")
                               ? backendEndpoint[1..]      // remove leading /
                               : backendEndpoint;

        var request = new HttpRequestMessage
        {
            Method = new HttpMethod(context.Request.Method),
            RequestUri = new Uri(client.BaseAddress + cleanBackendPath + query)
        };

        // *** copy the body ***
        request.Content = new StreamContent(context.Request.Body);

        // *** copy headers ***
        foreach (var h in context.Request.Headers)
            if (!request.Headers.TryAddWithoutValidation(h.Key, h.Value.ToArray()))
                request.Content?.Headers.TryAddWithoutValidation(h.Key, h.Value.ToArray());

        // log.LogInformation("Request URI={RequestURI}", request.RequestUri);

        using var response = await client.SendAsync(
                                request,
                                HttpCompletionOption.ResponseHeadersRead,
                                context.RequestAborted);

        context.Response.StatusCode = (int)response.StatusCode;
        await response.Content.CopyToAsync(context.Response.Body);
    }
});
if (app.Environment.IsDevelopment())
{
    app.UseSwagger();
    app.UseSwaggerUI();
}

app.UseHttpsRedirection();

app.MapGet("/", () => "Hello, World!");
app.Run();