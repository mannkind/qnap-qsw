using System.Diagnostics.CodeAnalysis;
using CommandLine;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Logging;
using QNAP.Commands;
using QNAP.DataAccess;

public class Program
{
    // These exist for AOT compilation
    [DynamicDependency(DynamicallyAccessedMemberTypes.All, typeof(GlobalOpts))]
    [DynamicDependency(DynamicallyAccessedMemberTypes.All, typeof(LoginOpts))]
    [DynamicDependency(DynamicallyAccessedMemberTypes.All, typeof(PoeModeOpts))]
    public static async Task<int> Main(string[] args)
    {
        // Create cancellation token
        var cts = new CancellationTokenSource();
        Console.CancelKeyPress += (s, e) => cts.Cancel();

        // Create DI service collection
        var services = new ServiceCollection();
        services.AddLogging((logging) =>
        {
            logging.AddFilter("System", LogLevel.Warning);
            logging.AddFilter("QNAP", LogLevel.Debug);
            logging.AddSimpleConsole(c =>
            {
                c.SingleLine = true;
                c.TimestampFormat = "[yyyy-MM-dd HH:mm:ss] ";
            });
        });
        services.AddSingleton<LoginCmd>();
        services.AddSingleton<PoeModeCmd>();
        services.AddSingleton<ISwitchDAO, SwitchDAO>();
        services.AddHttpClient(string.Empty).ConfigurePrimaryHttpMessageHandler(() => new HttpClientHandler
        {
            ClientCertificateOptions = ClientCertificateOption.Manual,
            ServerCertificateCustomValidationCallback =
                (httpRequestMessage, cert, cetChain, policyErrors) =>
                {
                    return true;
                }
        });

        // Build service provider
        var provider = services.BuildServiceProvider();

        // Parse command line arguments and run the appropriate command
        return await Parser
            .Default
            .ParseArguments<LoginOpts, PoeModeOpts>(args)
            .MapResult(
                async (LoginOpts opts) =>
                {
                    var cmd = provider.GetRequiredService<LoginCmd>();
                    return await cmd.RunAsync(opts, cts.Token);
                },
                async (PoeModeOpts opts) =>
                {
                    var cmd = provider.GetRequiredService<PoeModeCmd>();
                    return await cmd.RunAsync(opts, cts.Token);
                },
                _ => Task.FromResult(1)
            );
    }
}
