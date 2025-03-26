using BenchmarkDotNet.Attributes;
using BenchmarkDotNet.Running;
using Microsoft.Extensions.Caching.Memory;
using Microsoft.Extensions.DependencyInjection;


var summary = BenchmarkRunner.Run<MememoryCacheKey>();

[MemoryDiagnoser]
public class MememoryCacheKey
{
    private IMemoryCache? _cache;

    public IEnumerable<int> _iteration;

    public IEnumerable<string> _source;

    public MememoryCacheKey()
    {
        var provider = new ServiceCollection()
                               .AddMemoryCache()
                               .BuildServiceProvider();

        _cache = provider.GetService<IMemoryCache>();

        _source = Enumerable.Range(0, 100).Select(_ => DateTimeOffset.UtcNow.Ticks.ToString());

        _iteration = Enumerable.Range(0, 100).ToArray();
    }


    [Benchmark]
    public void TestStringKey()
    {
        foreach (var input in _source)
            foreach (var i in _iteration)
                _cache?.TryGetValue($"{input}_blob", out object? value);
    }

    [Benchmark]
    public void TestTupleKey()
    {
        foreach (var input in _source)
            foreach (var i in _iteration)
                _cache?.TryGetValue((input, "blob"), out object? value);
    }
}


