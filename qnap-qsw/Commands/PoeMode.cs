
using CommandLine;
using Microsoft.Extensions.Logging;
namespace QNAP.Commands;

public class PoeModeCmd : GlobalCmd
{
    public PoeModeCmd(ILogger<PoeModeCmd> logger, QNAP.DataAccess.ISwitchDAO dao) : base(logger, dao)
    {
    }

    /// <summary>
    /// 
    /// </summary>
    /// <param name="opts"></param>
    /// <param name="cancellationToken"></param>
    /// <returns></returns>
    public async Task<int> RunAsync(PoeModeOpts opts, CancellationToken cancellationToken = default)
    {
        this.Logger.LogDebug("Configure switch; host={host}, useHttps={useHttps}", opts.Host, opts.UseHttps);
        this.DAO.Configure(opts.UseHttps, opts.Host);

        this.Logger.LogDebug("Login to switch");
        var token = await this.DAO.LoginAsync(opts.Password, cancellationToken);
        if (string.IsNullOrEmpty(token))
        {
            this.Logger.LogError("Login to switch failed");
            return 1;
        }

        var interfaces = await this.DAO.POEInterfaces(cancellationToken);
        var portModes = new List<PortMode>{
            new() { Ports = opts.DisablePorts, Mode =  "disable" },
            new() { Ports = opts.PoePorts, Mode =  "poe" },
            new() { Ports = opts.PoePlusPorts, Mode =  "poe+" },
            new() { Ports = opts.PoePlusPlusPorts, Mode =  "poe++" },
        };
        var tasks = new List<Task<bool>>();
        foreach (var portMode in portModes)
        {
            // Skip empty lists of ports
            if (!portMode.Ports.Any(x => !string.IsNullOrWhiteSpace(x)))
            {
                continue;
            }

            foreach (var port in portMode.Ports)
            {
                this.Logger.LogDebug("Modifying poeMode; port={port}, mode={mode}", port, portMode.Mode);

                var properties = interfaces[port];
                var modifiedProperties = properties with { Mode = portMode.Mode };
                tasks.Add(this.DAO.UpdatePOEInterface(port, modifiedProperties, cancellationToken));
            }
        }

        this.Logger.LogDebug("Waiting for all poeMode modifications to complete");
        await Task.WhenAll(tasks);
        this.Logger.LogDebug("All poeMode modifications to complete");

        return 0;
    }
}

[Verb("poeMode", HelpText = "Change the POE port mode.")]
public record PoeModeOpts : GlobalOpts
{
    [Option("disable-ports", Required = false, Separator = ',', HelpText = "The ports to disable.")]
    public IEnumerable<string> DisablePorts { get; set; } = Array.Empty<string>();

    [Option("poe-ports", Required = false, Separator = ',', HelpText = "The ports to set to poe.")]
    public IEnumerable<string> PoePorts { get; set; } = Array.Empty<string>();

    [Option("poeplus-ports", Required = false, Separator = ',', HelpText = "The ports to set to poe+.")]
    public IEnumerable<string> PoePlusPorts { get; set; } = Array.Empty<string>();

    [Option("poeplusplus-ports", Required = false, Separator = ',', HelpText = "The ports to set to poe++.")]
    public IEnumerable<string> PoePlusPlusPorts { get; set; } = Array.Empty<string>();
}

record PortMode
{
    public IEnumerable<string> Ports { get; init; } = Array.Empty<string>();
    public string Mode { get; init; } = string.Empty;
}
