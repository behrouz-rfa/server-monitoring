package plugin

import (
	"server-monitoring/shared/view"

	"html/template"
	"log"
)

// TagHelper returns a template.FuncMap
// * JS returns JavaScript tag
// * CSS returns stylesheet tag
// * LINK returns hyperlink tag
func TagHelper(v view.View) template.FuncMap {
	f := make(template.FuncMap)

	f["JS"] = func(s string) template.HTML {
		path, err := v.AssetTimePath(s)

		if err != nil {
			log.Println("JS Error:", err)
			return template.HTML("")
		}

		return template.HTML(`<script type="text/javascript" src="` + path + `"></script>`)
	}

	f["CSS"] = func(s string) template.HTML {
		path, err := v.AssetTimePath(s)

		if err != nil {
			log.Println("CSS Error:", err)
			return template.HTML("<!-- CSS Error: " + s + " -->")
		}

		return template.HTML(`<link rel="stylesheet" type="text/css" href="` + path + `" />`)
	}

	f["LINK"] = func(path, name string) template.HTML {
		return template.HTML(`<a href="` + v.PrependBaseURI(path) + `">` + name + `</a>`)

	}
//
//	f["form"] = func(form elements.Form) template.HTML {
//		var html = ``
//		if len(form.Title) > 0 {
//			html += ` <h3 class="main_question"><strong>2/5</strong>` + form.Title + `</h3>`
//		}
//		if len(form.CheckBoxs) > 0 {
//			for _, element := range form.CheckBoxs {
//				html += `<label class="container_check">` + element.Title + `
//                                        <input type="checkbox" name="` + element.Name + `" value="` + element.Value + `" class="required">
//                                        <span class="checkmark"></span>
//                                    </label>`
//			}
//		}
//
//		return template.HTML(html)
//	}
//
//	f["element"] = func(element elements.TextInput) template.HTML {
//
//		switch element.Type {
//		case "text":
//			return template.HTML(`<div class="form-group">
//                                    <input type="text" name="` + element.Name + `" class="form-control" placeholder="` + element.Placeholder + `">
//                                </div>`)
//		case "password":
//			return template.HTML(`<div class="form-group">
//                                    <input type="password" name="` + element.Name + `" class="form-control" placeholder="` + element.Placeholder + `">
//                                </div>`)
//
//		case "email":
//			return template.HTML(`<div class="form-group">
//                                    <input type="email" name="` + element.Name + `" class="form-control" placeholder="` + element.Placeholder + `">
//                                </div>`)
//		}
//
//		return template.HTML(`<p>not found</p>`)
//	}
//	f["radio"] = func(radios  []elements.RadioButton) template.HTML {
//		html := ``
//
//		for _, element := range radios {
//			html += ` <div class="form-group radio_input"> <label class="container_radio">` + element.Title + `
//                                                <input  type="radio" name="` + element.Name + `" value="` + fmt.Sprintf("%d",element.Value) + `" >
//                                                <span class="checkmark"></span>
//                                            </label></div>`
//
//		}
//		html += ` `
//		return template.HTML(html)
//	}
//	f["checkbox"] = func(checkboxes  []elements.CheckBox) template.HTML {
//		html:=``
//		for _, element := range checkboxes {
//			html += `<label class="container_check">` + element.Title + `
//                                        <input type="checkbox" name="` + element.Name + `" value="` + element.Value + `" class="required">
//                                        <span class="checkmark"></span>
//                                    </label>`
//		}
//
//
//		return template.HTML(html)
//	}
//
//	f["textarea"] = func(name string, header string, radioElements interface{}) template.HTML {
//		return template.HTML(`
//
//<div class="form-group">
//                                    <label>` + header + `</label>
//                                    <textarea name="additional_message" class="form-control" style="height:100px;" placeholder="Type here additional info..." onkeyup="getVals(this, &#39;additional_message&#39;);"></textarea>
//                                </div>`)
//	}
//	f["select"] = func(seletes  []elements.SelectInput) template.HTML {
//		html := ``
//		for _, element := range seletes {
//			html += ` <div class="form-group">
//                                    <div class="styled-select clearfix">
//                                        <select class="wide required" name="` + element.Name + `" style="display: none;">`
//
//			for key, value := range element.Options {
//				html += `<option value="` + key + `">` + value + `</option>`
//			}
//			html += `</select><div class="nice-select wide required" tabindex="0"><span class="current">Your Country</span><ul class="list"><li data-value="" class="option selected">Your Country</li><li data-value="Europe" class="option">Europe</li><li data-value="Asia" class="option">Asia</li><li data-value="North America" class="option">North America</li><li data-value="South America" class="option">South America</li><li data-value="Oceania" class="option">Oceania</li></ul></div>
//                                    </div>
//                                </div>`
//
//		}
//
//		return template.HTML(html)
//	}
	return f
}
