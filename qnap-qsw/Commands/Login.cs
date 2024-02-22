
using CommandLine;
using Microsoft.Extensions.Logging;
namespace QNAP.Commands;

public class LoginCmd : GlobalCmd
{
    public LoginCmd(ILogger<LoginCmd> logger, QNAP.DataAccess.ISwitchDAO dao) : base(logger, dao)
    {
    }

    /// <summary>
    /// 
    /// </summary>
    /// <param name="opts"></param>
    /// <param name="cancellationToken"></param>
    /// <returns></returns>
    public async Task<int> RunAsync(LoginOpts opts, CancellationToken cancellationToken = default)
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

        Console.WriteLine(token);
        return 0;
    }
}

[Verb("login", HelpText = "Login to the QNAP QSW switch.")]
public record LoginOpts : GlobalOpts
{
}
