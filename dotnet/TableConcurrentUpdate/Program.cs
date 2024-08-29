// See https://aka.ms/new-console-template for more information


using Azure.Data.Tables;
using Azure.Identity;

try
{
    var storageUri = new Uri("https://test.table.core.windows.net/");
    var tableName = "test";
    var tableClient = new TableClient(storageUri, tableName, new DefaultAzureCredential());


    var partitionKey = "testpartitionkey";
    var rowkey = DateTimeOffset.Now.ToString("O");
    var propertyName = "Count";
    var entity = new TableEntity(partitionKey, rowkey)
    {
        {propertyName, 1}
    };

    var response = await tableClient.UpsertEntityAsync<TableEntity>(entity);

    Console.WriteLine("Test update without tag");

    Console.WriteLine(response.Content.ToString());

    var createdEntity = (await tableClient.GetEntityAsync<TableEntity>(partitionKey, rowkey)).Value;

    Console.WriteLine($"Initialized to '{createdEntity[propertyName]}'");

    createdEntity[propertyName] = 2;

    Console.WriteLine($"Updating to '{createdEntity[propertyName]}'");

    tableClient.SubmitTransaction([new TableTransactionAction(TableTransactionActionType.UpdateReplace, createdEntity)]);

    var updatedEntity = (await tableClient.GetEntityAsync<TableEntity>(partitionKey, rowkey)).Value;

    Console.WriteLine($"Updated to '{updatedEntity[propertyName]}'");

    createdEntity[propertyName] = 3;

    Console.WriteLine($"Updating to '{createdEntity[propertyName]}'");

    tableClient.SubmitTransaction([new TableTransactionAction(TableTransactionActionType.UpdateReplace, createdEntity)]);

    updatedEntity = (await tableClient.GetEntityAsync<TableEntity>(partitionKey, rowkey)).Value;

    Console.WriteLine($"Updated to '{updatedEntity[propertyName]}'");

    Console.WriteLine("====================");

    Console.WriteLine("Test update with tag");

    createdEntity = updatedEntity;

    createdEntity[propertyName] = 4;

    Console.WriteLine($"Updating to '{createdEntity[propertyName]}'");

    tableClient.SubmitTransaction([new TableTransactionAction(TableTransactionActionType.UpdateReplace, createdEntity, createdEntity.ETag)]);

    updatedEntity = (await tableClient.GetEntityAsync<TableEntity>(partitionKey, rowkey)).Value;

    Console.WriteLine($"Updated to '{updatedEntity[propertyName]}'");

    createdEntity[propertyName] = 5;

    Console.WriteLine($"Updating to '{createdEntity[propertyName]}'");

    try
    {
        tableClient.SubmitTransaction([new TableTransactionAction(TableTransactionActionType.UpdateReplace, createdEntity, createdEntity.ETag)]);
    }
    catch (Exception ex)
    {
        Console.WriteLine($"Update failed: {ex}");
    }
}
catch (Exception ex)
{
    Console.WriteLine($"Program failed: {ex}");
}
