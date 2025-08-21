namespace Project.Models
{
public class Items
{
    public int Id { get; set; }
    public string DogId { get; set; }      // From API "id"
    public string Name { get; set; }
    public string Description { get; set; }
    public int LifeMin { get; set; }
    public int LifeMax { get; set; }
    public int MaleWeightMin { get; set; }
    public int MaleWeightMax { get; set; }
    public int FemaleWeightMin { get; set; }
    public int FemaleWeightMax { get; set; }
    public bool Hypoallergenic { get; set; }
}



}