using System.Net.Http.Json;
using Project.Data;
using Project.Models;

public class DogService
{
    private readonly HttpClient _http;
    private readonly MyAppContext _context;

    public DogService(HttpClient http, MyAppContext context)
    {
        _http = http;
        _context = context;
    }

    public async Task FetchAndSaveDogsAsync()
    {
        var response = await _http.GetFromJsonAsync<DogApiResponse>("https://dogapi.dog/api/v2/breeds");

        if (response?.Data == null) return;

        foreach (var dog in response.Data)
        {
            var item = new Items
            {
                DogId = dog.Id,
                Name = dog.Attributes.Name,
                Description = dog.Attributes.Description,
                LifeMin = dog.Attributes.Life.Min,
                LifeMax = dog.Attributes.Life.Max,
                MaleWeightMin = dog.Attributes.MaleWeight.Min,
                MaleWeightMax = dog.Attributes.MaleWeight.Max,
                FemaleWeightMin = dog.Attributes.FemaleWeight.Min,
                FemaleWeightMax = dog.Attributes.FemaleWeight.Max,
                Hypoallergenic = dog.Attributes.Hypoallergenic??false
            };

            // Check if already exists
            if (!_context.Items.Any(i => i.DogId == dog.Id))
                _context.Items.Add(item);
        }

        await _context.SaveChangesAsync();
    }
}

// DTOs to match API JSON
public class DogApiResponse
{
    public List<DogData> Data { get; set; } = new List<DogData>();
}

public class DogData
{
    public required string Id { get; set; } = string.Empty;
    public  DogAttributes Attributes { get; set; } = new DogAttributes();
}


public class DogAttributes
{
    public string? Name { get; set; }
    public string? Description { get; set; }
    public DogLife? Life { get; set; }
    public DogWeight? MaleWeight { get; set; }
    public DogWeight? FemaleWeight { get; set; }
    public bool? Hypoallergenic { get; set; }
}
public DogAttributes Attributes { get; set; } = new DogAttributes();


public class DogLife { public int Min { get; set; } public int Max { get; set; } }
public class DogWeight { public int Min { get; set; } public int Max { get; set; } }
