// See https://aka.ms/new-console-template for more information

using Azure;
using Azure.Storage.Blobs;
using Azure.Storage.Blobs.Models;
using System.Formats.Tar;
using System.Reflection;
using System.Diagnostics;

var blobClient = new BlobClient(new Uri("https://bindudemo.blob.core.windows.net/public/tar.tar"));

using var blobStream = await blobClient.OpenReadAsync();

var reader = new TarReader(blobStream);

var entry = await reader.GetNextEntryAsync(false);

var tasks = new List<Task>();

var stopWatch = Stopwatch.StartNew();

while(entry != null)
{
    Console.WriteLine("Name: " + entry.Name);
    Console.WriteLine("Format: " + entry.Format);
    Console.WriteLine("Mode: " + entry.Mode);
    Console.WriteLine("Length: " + entry.Length);
    Console.WriteLine("Checksum: " + entry.Checksum);

    if (entry.DataStream != null)
    {
        var startPos = entry.DataStream?.GetType()?.GetField("_startInSuperStream", BindingFlags.NonPublic | BindingFlags.Instance)?.GetValue(entry.DataStream);
        var fileName = entry.Name;
        if (startPos != null)
        {
            Console.WriteLine("StartPos: " + startPos);
            var task = Task.Run(async () => 
            {
                var content = await blobClient.DownloadStreamingAsync(new BlobDownloadOptions
                {
                    Range = new HttpRange((long)startPos, entry.Length)
                });
                using var destination = new FileStream(Path.Combine("D:", fileName), FileMode.Create);
                await content.Value.Content.CopyToAsync(destination);
            });

            // await task; // Uncomment this to run synchronously
            tasks.Add(task);
        }
    }
    Console.WriteLine("=========================");

    entry = await reader.GetNextEntryAsync(false);
}

await Task.WhenAll(tasks);

Console.WriteLine("Time(seconds):" + stopWatch.Elapsed.TotalSeconds);