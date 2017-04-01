package main

import(
    "io/ioutil"
    "net/http"
    "html/template"

    "gopkg.in/gin-gonic/gin.v1"
    "github.com/russross/blackfriday"
    "github.com/microcosm-cc/bluemonday"
)

// ----------
// Create Page Construct for the single pages of this wiki
// ----------
type Page struct {
    Title string
    Body []byte
}

// Save a new Page
func (p *Page) Save() error  {
    filename := "./pages/" + p.Title + ".md"
    return ioutil.WriteFile(filename, p.Body, 0600)
}

// Load a page
func LoadPage(title string) (*Page, error) {
    filename := "./pages/" + title + ".md"
    body, err := ioutil.ReadFile(filename)
    if err != nil {
        return nil, err
    }
    return &Page{Title: title, Body: body}, nil
}

// ----------
// Utily funcs
// ----------


// ----------
// Handler for different routes
// ----------
func ViewHandler(c *gin.Context)  {
    title := c.Param("title")
    p, err := LoadPage(title)
    if err != nil {
      c.Redirect(http.StatusFound, "/edit/" + title)
    }
    p.Body = bluemonday.UGCPolicy().SanitizeBytes(blackfriday.MarkdownBasic(p.Body))
    c.HTML(http.StatusOK, "view.tmpl", gin.H{"Title": p.Title, "Body": template.HTML(p.Body)})
}
func EditHandler(c *gin.Context)  {
    title := c.Param("title")
    p, err := LoadPage(title)
    if err != nil {
        p = &Page{Title: title}
    }
    c.HTML(http.StatusOK, "edit.tmpl", gin.H{"Title": p.Title, "Body": p.Body})
}

func SaveHandler(c *gin.Context)  {
    title := c.Param("title")
    body := c.PostForm("body")
    p := &Page{Title: title, Body: []byte(body)}
    err := p.Save()
    if err != nil {
      // TODO: throw internal server error
      return
    }
    c.Redirect(http.StatusFound, "/view/" + title)
}

func HomeHandler(c *gin.Context) {
    c.Redirect(http.StatusFound, "/view/index")
}

func UploadHandler(c *gin.Context)  {

}

func main() {
    router := gin.Default()
    router.LoadHTMLGlob("./templates/*")
    router.Static("/assets", "./assets")

    html := template.Must(template.ParseFiles("./templates/edit.tmpl", "./templates/view.tmpl"))
    router.SetHTMLTemplate(html)

    router.GET("/", HomeHandler)
    router.GET("/view/:title", ViewHandler)
    router.GET("/edit/:title", EditHandler)
    router.POST("/save/:title", SaveHandler)
    router.POST("/upload/:title", UploadHandler)

    router.Run()
}
