package main

import (
	"encoding/json"
	"math"
	"os"

	"github.com/gin-gonic/gin"
)

type SinhVien struct {
	Masv string `json:"masv"`
	Ten  string `json:"ten"`
	Tuoi int    `json:"tuoi"`
}

type MonHoc struct {
	Mamon  string `json:"mamon"`
	Ten    string `json:"ten"`
	Tinchi int    `json:"tinchi"`
}

type DangKy struct {
	Masv  string  `json:"masv"`
	Mamon string  `json:"mamon"`
	Diem  float64 `json:"diem"`
}

type DiemDanh struct {
	Masv  string `json:"masv"`
	Mamon string `json:"mamon"`
	Buoi  int    `json:"buoi"`
	CoMat bool   `json:"comat"`
}

type LMS struct {
	SinhViens []SinhVien `json:"sinhviens"`
	MonHocs   []MonHoc   `json:"monhocs"`
	DangKys   []DangKy   `json:"dangkys"`
	DiemDanhs []DiemDanh `json:"diemdanhs"`
}

var data LMS

func save() {
	file, _ := json.MarshalIndent(data, "", "  ")
	os.WriteFile("data.json", file, 0644)
}

func load() {
	file, err := os.ReadFile("data.json")
	if err == nil {
		json.Unmarshal(file, &data)
	}
}

func tinhGPA(masv string) float64 {
	var tong float64
	var tin int

	for _, dk := range data.DangKys {
		if dk.Masv == masv {
			for _, m := range data.MonHocs {
				if m.Mamon == dk.Mamon {

					diem4 := diemHe4(dk.Diem)

					tong += diem4 * float64(m.Tinchi)
					tin += m.Tinchi
				}
			}
		}
	}

	if tin == 0 {
		return 0
	}
	return math.Round((tong/float64(tin))*100) / 100
}

func diemHe4(d float64) float64 {
	if d >= 8.5 {
		return 4.0
	} else if d >= 8.0 {
		return 3.5
	} else if d >= 7.0 {
		return 3.0
	} else if d >= 6.5 {
		return 2.5
	} else if d >= 5.5 {
		return 2.0
	} else if d >= 5.0 {
		return 1.5
	} else if d >= 4.0 {
		return 1.0
	}
	return 0
}

func main() {
	load()
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"msg": "LMS API is running",
			"endpoints": gin.H{
				"GET /sinhvien":        "Xem danh sach sinh vien",
				"POST /sinhvien":       "Them sinh vien",
				"DELETE /sinhvien/:id": "Xoa sinh vien",
				"GET /mon":             "Xem mon hoc",
				"POST /mon":            "Them mon hoc",
				"POST /dangky":         "Dang ky mon hoc",
				"DELETE /dangky":       "Huy dang ky mon hoc",
				"PUT /dangky":          "Xem mon da dang ky",
				"GET /gpa/:masv":       "Xem GPA",
				"POST /diemdanh":       "Điem danh",
				"GET /diemdanh/:mamon": "Xem danh sach diem danh theo mon",
				"DELETE /diemdanh":     "Xoa diem danh",
			},
			"how_to_test": "Su dung Postman hoac curl de kiem tra API",
		})
	})

	r.GET("/ui", func(c *gin.Context) {
		c.Header("Content-Type", "text/html")
		c.String(200, `
<!DOCTYPE html>
<html>
<head>
	<title>LMS Full UI</title>
	<style>
		body { font-family: Arial; padding: 20px; }
		input, button { margin: 5px; padding: 5px; }
		table { border-collapse: collapse; margin-top: 10px; }
		td, th { border: 1px solid black; padding: 6px; }
		h2 { margin-top: 30px; }
	</style>
</head>
<body>

<h1>LMS Management</h1>

<!-- ===== Sinh Vien ===== -->
<h2>Sinh Vien</h2>
<input id="sv_masv" placeholder="Ma SV">
<input id="sv_ten" placeholder="Ten">
<input id="sv_tuoi" placeholder="Tuoi">
<button onclick="themSV()">Them</button>
<button onclick="loadSV()">Load</button>

<table id="sv_table"></table>

<!-- ===== Mon hoc ===== -->
<h2>Mon hoc</h2>
<input id="mon_id" placeholder="Ma mon">
<input id="mon_ten" placeholder="Ten mon">
<input id="mon_tc" placeholder="Tin chi">
<button onclick="themMon()">Them</button>
<button onclick="loadMon()">Load</button>

<table id="mon_table"></table>

<!-- ===== Dang ky ===== -->
<h2>Dang Ky</h2>
<input id="dk_masv" placeholder="Ma SV">
<input id="dk_mamon" placeholder="Ma mon">
<button onclick="dangKy()">Dang Ky</button>

<h3>Nhap diem</h3>
<input id="dk_masv2" placeholder="Ma SV">
<input id="dk_mamon2" placeholder="Ma mon">
<input id="dk_diem" placeholder="Diem">
<button onclick="nhapDiem()">Cap nhat</button>

<!-- ===== GPA ===== -->
<h2>GPA</h2>
<input id="gpa_id" placeholder="Ma SV">
<button onclick="xemGPA()">Xem</button>
<p id="gpa_result"></p>

<!-- ===== DIEM DANH ===== -->
<h2>Diem danh</h2>
<input id="dd_masv" placeholder="Ma SV">
<input id="dd_mamon" placeholder="Ma môn">
<input id="dd_buoi" placeholder="Buoi">
<select id="dd_comat">
	<option value="true">Co mat</option>
	<option value="false">Vang</option>
</select>
<button onclick="themDD()">Them</button>

<h3>Xem diem danh theo mon</h3>
<input id="dd_mamon_xem" placeholder="Ma mon">
<button onclick="xemDD()">Xem</button>

<table id="dd_table"></table>

<script>

function themSV() {
	fetch('/sinhvien', {
		method: 'POST',
		headers: {'Content-Type': 'application/json'},
		body: JSON.stringify({
			masv: sv_masv.value,
			ten: sv_ten.value,
			tuoi: parseInt(sv_tuoi.value)
		})
	})
	.then(() => {
		alert("Them sinh vien thanh cong!");
		loadSV();
	});
}

function loadSV() {
	fetch('/sinhvien')
	.then(r => r.json())
	.then(data => {
		let html = "<tr><th>Ma</th><th>Ten</th><th>Tuoi</th><th>Xoa</th></tr>";
		data.forEach(sv => {
			html += "<tr><td>"+sv.masv+"</td><td>"+sv.ten+"</td><td>"+sv.tuoi+"</td>"
			+ "<td><button onclick=\"xoaSV('"+sv.masv+"')\">X</button></td></tr>";
		});
		sv_table.innerHTML = html;
	});
}

function xoaSV(id) {
	fetch('/sinhvien/'+id, {method:'DELETE'})
	.then(() => {
		alert("Da xoa!");
		loadSV();
	});
}

function themMon() {
	fetch('/mon', {
		method:'POST',
		headers:{'Content-Type':'application/json'},
		body:JSON.stringify({
			mamon:mon_id.value,
			ten:mon_ten.value,
			tinchi:parseInt(mon_tc.value)
		})
	})
	.then(() => {
		alert("Them mon thanh cong!");
		loadMon();
	});
}

function loadMon() {
	fetch('/mon')
	.then(r=>r.json())
	.then(data=>{
		let html="<tr><th>Ma</th><th>Ten</th><th>TC</th></tr>";
		data.forEach(m=>{
			html+="<tr><td>"+m.mamon+"</td><td>"+m.ten+"</td><td>"+m.tinchi+"</td></tr>";
		});
		mon_table.innerHTML=html;
	});
}

function dangKy() {
	fetch('/dangky', {
		method:'POST',
		headers:{'Content-Type':'application/json'},
		body:JSON.stringify({
			masv:dk_masv.value,
			mamon:dk_mamon.value
		})
	})
	.then(res => res.json())
	.then(() => {
		alert("Dang ky thanh cong!");
	});
}

function nhapDiem() {
	fetch('/dangky', {
		method:'PUT',
		headers:{'Content-Type':'application/json'},
		body:JSON.stringify({
			masv:dk_masv2.value,
			mamon:dk_mamon2.value,
			diem:parseFloat(dk_diem.value)
		})
	})
	.then(res => res.json())
	.then(() => {
		alert("Cap nhat diem thanh cong!");
	});
}

function xemGPA() {
	fetch('/gpa/'+gpa_id.value)
	.then(r=>r.json())
	.then(d=>{
		gpa_result.innerText="GPA: "+d.gpa;
	});
}

function themDD() {
	fetch('/diemdanh', {
		method:'POST',
		headers:{'Content-Type':'application/json'},
		body:JSON.stringify({
			masv:dd_masv.value,
			mamon:dd_mamon.value,
			buoi:parseInt(dd_buoi.value),
			comat:dd_comat.value==="true"
		})
	})
	.then(() => {
		alert("Diem danh thanh cong!");
	});
}

function xemDD() {
	fetch('/diemdanh/'+dd_mamon_xem.value)
	.then(r=>r.json())
	.then(data=>{
		let html="<tr><th>SV</th><th>Buoi</th><th>Co mat</th></tr>";
		data.forEach(d=>{
			html+="<tr><td>"+d.masv+"</td><td>"+d.buoi+"</td><td>"+d.comat+"</td></tr>";
		});
		dd_table.innerHTML=html;
	});
}

</script>

</body>
</html>
`)
	})

	r.POST("/sinhvien", func(c *gin.Context) {
		var sv SinhVien
		c.BindJSON(&sv)

		data.SinhViens = append(data.SinhViens, sv)
		save()

		c.JSON(200, sv)
	})

	r.GET("/sinhvien", func(c *gin.Context) {
		c.JSON(200, data.SinhViens)
	})

	r.DELETE("/sinhvien/:id", func(c *gin.Context) {
		id := c.Param("id")

		for i, sv := range data.SinhViens {
			if sv.Masv == id {
				data.SinhViens = append(data.SinhViens[:i], data.SinhViens[i+1:]...)
				save()
				c.JSON(200, gin.H{"msg": "deleted"})
				return
			}
		}

		c.JSON(404, gin.H{"error": "not found"})
	})

	r.POST("/mon", func(c *gin.Context) {
		var m MonHoc
		c.BindJSON(&m)

		data.MonHocs = append(data.MonHocs, m)
		save()

		c.JSON(200, m)
	})

	r.GET("/mon", func(c *gin.Context) {
		c.JSON(200, data.MonHocs)
	})

	r.POST("/dangky", func(c *gin.Context) {
		var dk DangKy
		c.BindJSON(&dk)

		data.DangKys = append(data.DangKys, dk)
		save()

		c.JSON(200, dk)
	})

	r.PUT("/dangky", func(c *gin.Context) {
		var dk DangKy
		c.BindJSON(&dk)

		for i := range data.DangKys {
			if data.DangKys[i].Masv == dk.Masv && data.DangKys[i].Mamon == dk.Mamon {
				data.DangKys[i].Diem = dk.Diem
				save()
				c.JSON(200, gin.H{"msg": "updated"})
				return
			}
		}

		c.JSON(404, gin.H{"error": "not found"})
	})

	r.DELETE("/dangky", func(c *gin.Context) {
		var dk DangKy
		c.BindJSON(&dk)

		for i := range data.DangKys {
			if data.DangKys[i].Masv == dk.Masv && data.DangKys[i].Mamon == dk.Mamon {
				data.DangKys = append(data.DangKys[:i], data.DangKys[i+1:]...)
				save()
				c.JSON(200, gin.H{"msg": "deleted"})
				return
			}
		}

		c.JSON(404, gin.H{"error": "not found"})
	})

	r.GET("/gpa/:masv", func(c *gin.Context) {
		id := c.Param("masv")
		c.JSON(200, gin.H{"gpa": tinhGPA(id)})
	})

	r.POST("/diemdanh", func(c *gin.Context) {
		var dd DiemDanh
		c.BindJSON(&dd)

		data.DiemDanhs = append(data.DiemDanhs, dd)
		save()

		c.JSON(200, dd)
	})

	r.GET("/diemdanh/:mamon", func(c *gin.Context) {
		mamon := c.Param("mamon")

		var result []DiemDanh
		for _, d := range data.DiemDanhs {
			if d.Mamon == mamon {
				result = append(result, d)
			}
		}

		c.JSON(200, result)
	})

	r.DELETE("/diemdanh", func(c *gin.Context) {
		var dd DiemDanh
		c.BindJSON(&dd)

		for i, d := range data.DiemDanhs {
			if d.Masv == dd.Masv && d.Mamon == dd.Mamon && d.Buoi == dd.Buoi {
				data.DiemDanhs = append(data.DiemDanhs[:i], data.DiemDanhs[i+1:]...)
				save()
				c.JSON(200, gin.H{"msg": "deleted"})
				return
			}
		}

		c.JSON(404, gin.H{"error": "not found"})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r.Run(":" + port)
}
