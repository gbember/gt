// main.go
package main

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/gbember/gt/navmesh"
	"github.com/google/gxui"
	"github.com/google/gxui/drivers/gl"
	"github.com/google/gxui/math"
	"github.com/google/gxui/samples/flags"
)

var (
	nmj *navmesh.NavMeshJson
	nm  *navmesh.NavMesh
	nmastar = navmesh.NavmeshAstar()
)

func main() {
	log.SetFlags(log.Lshortfile | log.Ldate)

	meshFileName := "../mesh.json"

	data, err := ioutil.ReadFile(meshFileName)
	if err != nil {
		log.Fatal(err)
	}
	nmj = new(navmesh.NavMeshJson)
	err = json.Unmarshal(data, nmj)
	if err != nil {
		log.Fatal(err)
	}

	n_m, err := navmesh.NewNavMesh(meshFileName)
	if err != nil {
		log.Fatal(err)
	}
	nm = n_m

	gl.StartDriver(appMain)

	//
	//	ps, isWalk := nm.FindPath(179, 41, 178, 886)
	//	log.Println(isWalk, ps)
	//	if isWalk {
	//		fn := "tt.cpuprof"
	//		f, err := os.Create(fn)
	//		if err != nil {
	//			log.Fatal(err)
	//		}
	//		err = pprof.StartCPUProfile(f)
	//		if err != nil {
	//			log.Fatal(err)
	//		}
	//		max := int64(100000)
	//		st := time.Now()
	//		for i := int64(0); i < max; i++ {
	//			nm.FindPath(179, 41, 178, 886)
	//		}
	//		nt := time.Since(st)
	//		log.Println(nt, nt.Nanoseconds()/max)

	//		pprof.StopCPUProfile()
	//	}
}

func appMain(driver gxui.Driver) {
	theme := flags.CreateTheme(driver)
	width := int(nmj.Width)
	height := int(nmj.Heigth)

	window := theme.CreateWindow(width, height, "navmesh")
	canvas := driver.CreateCanvas(math.Size{W: width, H: height})

	ps := nmj.Points

	// mouse
	isStart := true
	x1, y1, x2, y2 := int64(0), int64(0), int64(0), int64(0)

	window.OnMouseDown(func(me gxui.MouseEvent) {
		if nm.IsWalkOfPoint(navmesh.Point{X: int64(me.Point.X), Y: int64(me.Point.Y)}) {
			if isStart {
				x1 = int64(me.Point.X)
				y1 = int64(me.Point.Y)
			} else {
				x2 = int64(me.Point.X)
				y2 = int64(me.Point.Y)
			}
			if !isStart {
				drawWalkPath(window, theme, driver, x1, y1, x2, y2)
			}
			isStart = !isStart
		}
	})

	// draw mesh
	for i := 0; i < len(ps); i++ {
		polys := make([]gxui.PolygonVertex, 0, len(ps[i]))
		for j := 0; j < len(ps[i]); j++ {
			polys = append(polys, gxui.PolygonVertex{
				Position: math.Point{
					int(ps[i][j].X),
					int(ps[i][j].Y),
				}})
		}
		//		canvas.DrawPolygon(polys, gxui.CreatePen(2, gxui.Gray80), gxui.CreateBrush(gxui.Gray40))
		canvas.DrawPolygon(polys, gxui.CreatePen(2, gxui.Red), gxui.CreateBrush(gxui.Yellow))
	}

	canvas.Complete()
	image := theme.CreateImage()
	image.SetCanvas(canvas)
	window.AddChild(image)
	window.OnClose(driver.Terminate)
}

//画行走路线
func drawWalkPath(window gxui.Window, theme gxui.Theme, driver gxui.Driver, x1, y1, x2, y2 int64) {
	ps, isWalk := nm.FindPath(nmastar,x1, y1, x2, y2)
	if !isWalk {
		return
	}
	canvas := driver.CreateCanvas(math.Size{W: int(nmj.Width), H: int(nmj.Heigth)})

	var polys []gxui.PolygonVertex
	for i := 0; i < len(ps); i++ {

		polys = append(polys,
			gxui.PolygonVertex{
				Position: math.Point{
					int(ps[i].X),
					int(ps[i].Y),
				}})
	}
	canvas.DrawLines(polys, gxui.CreatePen(2, gxui.Green))

	canvas.Complete()
	image := theme.CreateImage()
	image.SetCanvas(canvas)
	window.AddChild(image)

}
