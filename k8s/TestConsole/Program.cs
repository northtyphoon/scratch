using System;
using System.Net.Http;
using System.Security.Cryptography.X509Certificates;
using System.Threading.Tasks;

namespace TestConsole
{
    class Program
    {
        static async Task Main(string[] args)
        {
            if (args.Length != 2)
            {
                Console.WriteLine("dotnet TestConsole <client-cert-pfx-file-path> <serivce-endpoint>");
                return;
            }

            var clientCertPfxFilePath = args[0];
            var serivceEndpoint =  args[1];

            try
            {
                await TestGetAsync(serivceEndpoint, clientCertPfxFilePath).ConfigureAwait(false);
            }
            catch (Exception ex)
            {
                Console.WriteLine($"Failed: {ex}");
            }
        }

        private static async Task TestGetAsync(string serivceEndpoint, string clientCertPfxFilePath)
        {
            var cert = new X509Certificate2(clientCertPfxFilePath, string.Empty);
            var handler = new HttpClientHandler();
            handler.ClientCertificates.Add(cert);
            handler.ServerCertificateCustomValidationCallback = (httpRequestMessage, cert, cetChain, policyErrors) =>
            {
                return true;
            };
            var client = new HttpClient(handler);

            var request = new HttpRequestMessage()
            {
                RequestUri = new Uri(serivceEndpoint),
                Method = HttpMethod.Get,
            };
            var response = await client.SendAsync(request);

            response.EnsureSuccessStatusCode();

            var responseContent = await response.Content.ReadAsStringAsync();
            Console.WriteLine($"Response: {responseContent}");
        }
    }
}
