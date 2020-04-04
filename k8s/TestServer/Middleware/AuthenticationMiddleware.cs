using Microsoft.AspNetCore.Http;
using System;
using System.Threading.Tasks;

namespace TestServer.Middleware
{
    public class AuthenticationMiddleware
    {
        private readonly RequestDelegate next;

        public AuthenticationMiddleware(RequestDelegate next)
        {
            this.next = next;
        }

        public async Task Invoke(HttpContext context)
        {
            await GetClientCertificateThumbprintAsync(context).ConfigureAwait(false);
            await next(context).ConfigureAwait(false);
        }

        private async Task GetClientCertificateThumbprintAsync(HttpContext context)
        {
            Console.WriteLine("=== Begin GetClientCertificateThumbprintAsync ===");

            foreach (var header in context.Request.Headers)
            {
                Console.WriteLine($"{header.Key}: {header.Value}");
            }

            var clientCert = await context.Connection.GetClientCertificateAsync().ConfigureAwait(false);
            Console.WriteLine($"GetClientCertificateAsync: {clientCert?.Thumbprint}");

            Console.WriteLine("=== End GetClientCertificateThumbprintAsync ===");
        }
    }
}
