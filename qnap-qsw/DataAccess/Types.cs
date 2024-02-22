using System.Text.Json.Serialization;
namespace QNAP.DataAccess;

public record LoginRequest
{
    [property: JsonPropertyName("username")]
    public string Username { get; init; } = string.Empty;

    [property: JsonPropertyName("password")]
    public string Password { get; init; } = string.Empty;
}

public record UpdatePOEInterfaceRequest
{
    [property: JsonPropertyName("idx")]
    public string Idx { get; init; } = string.Empty;

    [property: JsonPropertyName("data")]
    public POEInterfacesResponseKVPairValue Data { get; init; } = new();
}

public record LoginResponse
{
    [property: JsonPropertyName("error_code")]
    public int ErrorCode { get; init; } = 0;

    [property: JsonPropertyName("error_message")]
    public string ErrorMessage { get; init; } = string.Empty;

    [property: JsonPropertyName("result")]
    public LoginResponseTokens Result { get; init; } = new();
}

public record LoginResponseTokens
{
    public string AccessToken { get; init; } = string.Empty;
    public string RefreshToken { get; init; } = string.Empty;
}

public record POEInterfacesResponse
{
    [property: JsonPropertyName("error_code")]
    public int ErrorCode { get; init; } = 0;

    [property: JsonPropertyName("error_message")]
    public string ErrorMessage { get; init; } = string.Empty;

    [property: JsonPropertyName("result")]
    public IEnumerable<POEInterfacesResponseKVPair> Result { get; init; } = new List<POEInterfacesResponseKVPair>();
}

public record POEInterfacesResponseKVPair
{
    [property: JsonPropertyName("key")]
    public string Key { get; init; } = string.Empty;

    [property: JsonPropertyName("val")]
    public POEInterfacesResponseKVPairValue Value { get; init; } = new();
}

public record POEInterfacesResponseKVPairValue
{
    [property: JsonPropertyName("Mode")]
    public string Mode { get; init; } = string.Empty;

    [property: JsonPropertyName("Priority")]
    public string Priority { get; init; } = string.Empty;

    [property: JsonPropertyName("MaxPower")]
    public int MaxPower { get; init; } = 0;
}
