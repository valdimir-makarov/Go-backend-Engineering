using System.Security.Cryptography.X509Certificates;
using Microsoft.EntityFrameworkCore;
using Project.Models;
namespace Project.Data
{
    public class MyAppContext : DbContext
    {
        public MyAppContext(DbContextOptions<MyAppContext> options) : base(options){}
     
            public DbSet<Items> Items { get; set; }
       
     }
}