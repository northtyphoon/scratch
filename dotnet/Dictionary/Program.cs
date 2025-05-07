// See https://aka.ms/new-console-template for more information
using System.Collections.Frozen;

Console.WriteLine("Hello, World!");


var dict = new Dictionary<string, int>(StringComparer.OrdinalIgnoreCase)
{
    ["tags"] = 1,
    ["manifests"] = 2,
    ["blobs"] = 3
};
var lookup = dict.GetAlternateLookup<ReadOnlySpan<char>>();

var dict2 = dict.ToFrozenDictionary(StringComparer.OrdinalIgnoreCase);
var lookup2 = dict2.GetAlternateLookup<ReadOnlySpan<char>>();

ReadOnlySpan<char> key1 = "Blobs";
ReadOnlySpan<char> key2 = "blobs";

Console.WriteLine(lookup.TryGetValue(key1, out int _));
Console.WriteLine(lookup.TryGetValue(key2, out int __));

Console.WriteLine(lookup2.TryGetValue(key1, out int ___));
Console.WriteLine(lookup2.TryGetValue(key2, out int _____));