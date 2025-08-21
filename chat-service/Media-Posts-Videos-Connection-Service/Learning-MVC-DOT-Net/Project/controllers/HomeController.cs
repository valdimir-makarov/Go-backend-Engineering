using Microsoft.AspNetCore.Mvc;
using Project.Models;
namespace Project.Controllers
{
    public class HomeController : Controller
    {
        public IActionResult OverView()
        {
            //itialize the items object
            //in the Models folder
            var items = new Items()
            {
                Id = 1,
                Name = "Sample Item"
            };

            return View(items);
        }
        
    
     public IActionResult GetId(int id)
        {
            return Content($"Item ID: {id}");
        }
 }
}