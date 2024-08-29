using Microsoft.Extensions.Caching.Memory;
using Microsoft.Extensions.DependencyInjection;

var provider = new ServiceCollection()
                       .AddMemoryCache()
                       .BuildServiceProvider();

var cache = provider.GetService<IMemoryCache>();

cache.Set("key1", "value1", new MemoryCacheEntryOptions
{
    Size = 1,
    AbsoluteExpiration = DateTime.UtcNow.AddSeconds(5)
});

await Task.Delay(1000);

var value = cache.Get("key1");

Console.WriteLine(value);