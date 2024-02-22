using System.Diagnostics.CodeAnalysis;
using System.Net.Http.Json;
using System.Text.Json;
using System.Text.Json.Serialization;
using Microsoft.Extensions.Logging;
namespace QNAP.DataAccess;

public interface ISwitchDAO
{
    void Configure(bool https, string host);
    Task<string> LoginAsync(string password, CancellationToken cancellationToken = default);
    Task<Dictionary<string, POEInterfacesResponseKVPairValue>> POEInterfaces(CancellationToken cancellationToken = default);
    Task<bool> UpdatePOEInterface(string port, POEInterfacesResponseKVPairValue properties, CancellationToken cancellationToken);
}

public class SwitchDAO : ISwitchDAO
{

    /// <summary>
    /// Initializes a new instance of the SourceDAO class.
    /// </summary>
    /// <param name="logger"></param>
    /// <param name="httpClientFactory"></param>
    /// <returns></returns>
    public SwitchDAO(ILogger<SwitchDAO> logger, IHttpClientFactory httpClientFactory)
    {
        this.Logger = logger;
        this.Client = httpClientFactory.CreateClient();
    }

    /// <summary>
    /// 
    /// </summary>
    /// <param name="useHttps"></param>
    /// <param name="host"></param>
    /// <param name="token"></param>
    public void Configure(bool useHttps, string host)
    {
        this.Logger.LogDebug("Configuring switch; host={host}, useHttps={useHttps}", host, useHttps);
        this.UseHttps = useHttps;
        this.Host = host;
    }

    /// <summary>
    /// 
    /// </summary>
    /// <param name="password"></param>
    /// <param name="cancellationToken"></param>
    /// <returns></returns>
    [UnconditionalSuppressMessage("Aot", "IL2026:RequiresUnreferencedCodeAttribute", Justification = "Using QNAPSerializerContext")]
    [UnconditionalSuppressMessage("Aot", "IL3050:RequiresDynamicCode", Justification = "Using QNAPSerializerContext")]
    public async Task<string> LoginAsync(string password, CancellationToken cancellationToken = default)
    {
        this.Logger.LogDebug("Sending request to login");
        var request = this.Request(HttpMethod.Post, "users/login", new LoginRequest
        {
            Username = "admin",
            Password = Convert.ToBase64String(System.Text.Encoding.UTF8.GetBytes(password))
        });
        var resp = await this.Client.SendAsync(request, cancellationToken);
        if (resp.StatusCode != System.Net.HttpStatusCode.OK)
        {
            this.Logger.LogDebug("Request to login failed; statusCode={statusCode}", resp.StatusCode);
            return string.Empty;
        }

        this.Logger.LogDebug("Request to login successful");
        var content = await resp.Content.ReadAsStringAsync();
        var obj = JsonSerializer.Deserialize<LoginResponse>(content, options: new JsonSerializerOptions() { TypeInfoResolver = QNAPSerializerContext.Default });
        this.Token = obj?.Result.AccessToken ?? string.Empty;

        return this.Token;
    }

    /// <summary>
    /// 
    /// </summary>
    /// <param name="cancellationToken"></param>
    /// <returns></returns>
    [UnconditionalSuppressMessage("Aot", "IL2026:RequiresUnreferencedCodeAttribute", Justification = "Using QNAPSerializerContext")]
    [UnconditionalSuppressMessage("Aot", "IL3050:RequiresDynamicCode", Justification = "Using QNAPSerializerContext")]
    public async Task<Dictionary<string, POEInterfacesResponseKVPairValue>> POEInterfaces(CancellationToken cancellationToken = default)
    {
        this.Logger.LogDebug("Sending request to obtain poe interface statuses");
        var result = new Dictionary<string, POEInterfacesResponseKVPairValue>();
        var request = this.Request(HttpMethod.Get, "poe/interface");
        var resp = await this.Client.SendAsync(request, cancellationToken);
        if (resp.StatusCode != System.Net.HttpStatusCode.OK)
        {
            this.Logger.LogError("Request to obtain poe interface statuses failed; statusCode={statusCode}", resp.StatusCode);
            return result;
        }

        this.Logger.LogDebug("Request to obtain poe interface statuses successful");
        var content = await resp.Content.ReadAsStringAsync();
        var obj = JsonSerializer.Deserialize<POEInterfacesResponse>(content, options: new JsonSerializerOptions() { TypeInfoResolver = QNAPSerializerContext.Default });
        foreach (var item in obj?.Result ?? new List<POEInterfacesResponseKVPair>())
        {
            result.Add(item.Key, item.Value);
        }

        return result;
    }

    /// <summary>
    /// 
    /// </summary>
    /// <param name="port"></param>
    /// <param name="properties"></param>
    /// <param name="cancellationToken"></param>
    /// <returns></returns>
    public async Task<bool> UpdatePOEInterface(string port, POEInterfacesResponseKVPairValue properties, CancellationToken cancellationToken)
    {
        this.Logger.LogDebug("Sending request to set poe interface mode; port={port}, mode={mode}", port, properties.Mode);
        var request = this.Request(HttpMethod.Put, "poe/interface", new UpdatePOEInterfaceRequest
        {
            Idx = port,
            Data = properties with { Mode = this.modeTranslation.GetValueOrDefault(properties.Mode) ?? "unknown" },
        });
        var resp = await this.Client.SendAsync(request, cancellationToken);
        if (resp.StatusCode != System.Net.HttpStatusCode.OK)
        {
            this.Logger.LogError("Request to set poe interface mode failed; port={port}, mode={mode}", port, properties.Mode);
            return false;
        }

        this.Logger.LogDebug("Request to set poe interface mode sucessful; port={port}, mode={mode}", port, properties.Mode);
        return true;
    }

    /// <summary>
    /// 
    /// </summary>
    /// <param name="method"></param>
    /// <param name="path"></param>
    /// <param name="data"></param>
    /// <returns></returns>
    [UnconditionalSuppressMessage("Aot", "IL2026:RequiresUnreferencedCodeAttribute", Justification = "Using QNAPSerializerContext")]
    [UnconditionalSuppressMessage("Aot", "IL3050:RequiresDynamicCode", Justification = "Using QNAPSerializerContext")]
    private HttpRequestMessage Request(HttpMethod method, string path, object? data = null)
    {
        var protocol = this.UseHttps ? "https" : "http";
        var url = $"{protocol}://{this.Host}/api/v3/{path}";
        var request = new HttpRequestMessage(method, url);
        request.Headers.TryAddWithoutValidation("Content-Type", "application/json");

        // Add the bearer token if it exists
        if (!string.IsNullOrEmpty(this.Token))
        {
            request.Headers.TryAddWithoutValidation("Authorization", $"Bearer {this.Token}");
        }

        // Encode the json if it exits
        if (data != null)
        {
            request.Content = JsonContent.Create(data, options: new JsonSerializerOptions() { TypeInfoResolver = QNAPSerializerContext.Default });
        }

        return request;
    }

    /// <summary>
    /// The logger to use for the API
    /// </summary>
    private readonly ILogger<SwitchDAO> Logger;

    /// <summary>
    /// The client to use for the API
    /// </summary>
    private readonly HttpClient Client;

    /// <summary>
    /// The flag that indicates to use https or not
    /// </summary>
    private bool UseHttps { get; set; } = true;

    /// <summary>
    /// The host to use for the API
    /// </summary>
    private string Host { get; set; } = string.Empty;

    /// <summary>
    /// The token to use for the API
    /// </summary>
    private string Token { get; set; } = string.Empty;

    /// <summary>
    /// The translation of the mode for the API
    /// </summary>
    private Dictionary<string, string> modeTranslation = new Dictionary<string, string>
    {
        { "disable", "disable" },
        { "poe", "poeDot3af" },
        { "poe+", "poePlusDot3at" },
        { "poe++", "poePlusDot3bt" },
    };
}

[JsonSerializable(typeof(LoginRequest))]
[JsonSerializable(typeof(LoginResponse))]
[JsonSerializable(typeof(LoginResponseTokens))]
[JsonSerializable(typeof(POEInterfacesResponse))]
[JsonSerializable(typeof(POEInterfacesResponseKVPair))]
[JsonSerializable(typeof(POEInterfacesResponseKVPairValue))]
[JsonSerializable(typeof(UpdatePOEInterfaceRequest))]
[JsonSerializable(typeof(string))]
[JsonSerializable(typeof(int))]
public partial class QNAPSerializerContext : JsonSerializerContext
{
}
