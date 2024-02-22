using CommandLine;
using Microsoft.Extensions.Logging;
namespace QNAP.Commands;

public abstract class GlobalCmd
{
    public GlobalCmd(ILogger<GlobalCmd> logger, QNAP.DataAccess.ISwitchDAO dao)
    {
        this.Logger = logger;
        this.DAO = dao;
    }

    /// <summary>
    /// 
    /// </summary>
    protected ILogger<GlobalCmd> Logger;

    /// <summary>
    /// 
    /// </summary>
    protected QNAP.DataAccess.ISwitchDAO DAO;
}

public record GlobalOpts
{
    [Option("host", Required = false, HelpText = "The host/ip.")]
    public string Host { get; set; } = "switch.lan";

    [Option("password", Required = false, HelpText = "The password of the admin user (default: $QNAP_QSW_PASSWORD).")]
    public string Password { get; set; } = Environment.GetEnvironmentVariable("QNAP_QSW_PASSWORD") ?? "";

    [Option("use-https", Required = false, HelpText = "The flag that indicates if https should be used. (default: true)")]
    public bool UseHttps { get; set; } = true;
}
