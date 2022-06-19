Task CreateTask(CancellationToken cancellationToken)
{
    using var linkedCts = CancellationTokenSource.CreateLinkedTokenSource(cancellationToken);
    {
        return Task.Delay(5000, linkedCts.Token);
    }// Because the linkedCts is disposed, the cancellationToken won't be able to cancel the task.
}

var cts = new CancellationTokenSource();

var task = CreateTask(cts.Token);

cts.Cancel();

try
{
    await task.ConfigureAwait(false);

    Console.WriteLine($"IsCompleted: {task.IsCompleted}");
    Console.WriteLine($"IsFaulted: {task.IsFaulted}");
}
catch (Exception ex)
{
    Console.WriteLine(ex);
}

