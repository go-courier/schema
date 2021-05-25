/*
Package routes GENERATED BY gengo:validator 
DON'T EDIT THIS FILE
*/
package routes

import (
	mime_multipart "mime/multipart"

	github_com_go_courier_schema_pkg_validator "github.com/go-courier/schema/pkg/validator"
)

func validate_Cookie(v *Cookie) error {
	errSet := github_com_go_courier_schema_pkg_validator.NewErrorSet()
	if errSet.Len() > 0 {
		return errSet
	}
	return nil

}
func (v *Cookie) Validate() error {
	return validate_Cookie(v)
}
func validate_Create(v *Create) error {
	errSet := github_com_go_courier_schema_pkg_validator.NewErrorSet()
	errSet.AddErr(validate_Create_Field_Data(&v.Data), github_com_go_courier_schema_pkg_validator.Location("body"), "")
	if errSet.Len() > 0 {
		return errSet
	}
	return nil

}
func validate_Create_Field_Data(v *Data) error {
	return validate_Data(v)
}
func validate_Data(v *Data) error {
	errSet := github_com_go_courier_schema_pkg_validator.NewErrorSet()
	errSet.AddErr(validate_Data_Field_Data(&v.Data), "data")
	errSet.AddErr(validate_Data_Field_SubData(&v.SubData), "subData")
	if errSet.Len() > 0 {
		return errSet
	}
	return nil

}
func validate_Data_Field_Data(v **Data) error {
	vv := *v
	if vv == nil {
		return nil
	}
	return validate_Data_Field_Data_Ptr(vv)

}
func validate_Data_Field_Data_Ptr(v *Data) error {
	return validate_Data(v)
}
func validate_Data_Field_SubData(v **SubData) error {
	vv := *v
	if vv == nil {
		return nil
	}
	return validate_Data_Field_SubData_Ptr(vv)

}
func validate_Data_Field_SubData_Ptr(v *SubData) error {
	return validate_SubData(v)
}
func validate_SubData(v *SubData) error {
	errSet := github_com_go_courier_schema_pkg_validator.NewErrorSet()
	if errSet.Len() > 0 {
		return errSet
	}
	return nil

}
func validate_Create_Field_Data_Field_Data(v **Data) error {
	vv := *v
	if vv == nil {
		return nil
	}
	return validate_Create_Field_Data_Field_Data_Ptr(vv)

}
func validate_Create_Field_Data_Field_Data_Ptr(v *Data) error {
	return validate_Data(v)
}
func validate_Create_Field_Data_Field_SubData(v **SubData) error {
	vv := *v
	if vv == nil {
		return nil
	}
	return validate_Create_Field_Data_Field_SubData_Ptr(vv)

}
func validate_Create_Field_Data_Field_SubData_Ptr(v *SubData) error {
	return validate_SubData(v)
}
func (v *Create) Validate() error {
	return validate_Create(v)
}
func validate_DataProvider(v *DataProvider) error {
	errSet := github_com_go_courier_schema_pkg_validator.NewErrorSet()
	errSet.AddErr(validate_DataProvider_Field_ID(&v.ID), github_com_go_courier_schema_pkg_validator.Location("path"), "id")
	if errSet.Len() > 0 {
		return errSet
	}
	return nil

}
func validate_DataProvider_Field_ID(v *string) error {
	vv := *v
	if vv == "" {
		return &(github_com_go_courier_schema_pkg_validator.MissingRequired{})
	}
	return (&(github_com_go_courier_schema_pkg_validator.StringValidator{
		MinLength: 6,
	})).Validate(string(vv))

}
func (v *DataProvider) Validate() error {
	return validate_DataProvider(v)
}
func validate_DownloadFile(v *DownloadFile) error {
	errSet := github_com_go_courier_schema_pkg_validator.NewErrorSet()
	if errSet.Len() > 0 {
		return errSet
	}
	return nil

}
func (v *DownloadFile) Validate() error {
	return validate_DownloadFile(v)
}
func validate_FormMultipartWithFile(v *FormMultipartWithFile) error {
	errSet := github_com_go_courier_schema_pkg_validator.NewErrorSet()
	errSet.AddErr(validate_FormMultipartWithFile_Field_FormData(&v.FormData), github_com_go_courier_schema_pkg_validator.Location("body"), "")
	if errSet.Len() > 0 {
		return errSet
	}
	return nil

}
func validate_FormMultipartWithFile_Field_FormData(v *struct {
	String string                     `name:"string,omitempty"`
	Slice  []string                   `name:"slice,omitempty"`
	Data   Data                       `name:"data,omitempty"`
	File   *mime_multipart.FileHeader `name:"file" validate:"-"`
}) error {
	errSet := github_com_go_courier_schema_pkg_validator.NewErrorSet()
	errSet.AddErr(validate_FormMultipartWithFile_Field_FormData_Field_Slice(&v.Slice), "slice")
	errSet.AddErr(validate_FormMultipartWithFile_Field_FormData_Field_Data(&v.Data), "data")
	if errSet.Len() > 0 {
		return errSet
	}
	return nil

}
func validate_FormMultipartWithFile_Field_FormData_Field_Slice(v *[]string) error {
	list := *v
	n := len(list)
	errSet := github_com_go_courier_schema_pkg_validator.NewErrorSet()
	for i := 0; i < n; i++ {
		errSet.AddErr(validate_FormMultipartWithFile_Field_FormData_Field_Slice_Elem(&list[i]), i)
	}
	if errSet.Len() > 0 {
		return errSet
	}
	return nil

}
func validate_FormMultipartWithFile_Field_FormData_Field_Slice_Elem(v *string) error {
	return nil
}
func validate_FormMultipartWithFile_Field_FormData_Field_Data(v *Data) error {
	return validate_Data(v)
}
func (v *FormMultipartWithFile) Validate() error {
	return validate_FormMultipartWithFile(v)
}
func validate_FormMultipartWithFiles(v *FormMultipartWithFiles) error {
	errSet := github_com_go_courier_schema_pkg_validator.NewErrorSet()
	errSet.AddErr(validate_FormMultipartWithFiles_Field_FormData(&v.FormData), github_com_go_courier_schema_pkg_validator.Location("body"), "")
	if errSet.Len() > 0 {
		return errSet
	}
	return nil

}
func validate_FormMultipartWithFiles_Field_FormData(v *struct {
	Files []*mime_multipart.FileHeader `name:"files" validate:"-"`
}) error {
	errSet := github_com_go_courier_schema_pkg_validator.NewErrorSet()
	if errSet.Len() > 0 {
		return errSet
	}
	return nil

}
func (v *FormMultipartWithFiles) Validate() error {
	return validate_FormMultipartWithFiles(v)
}
func validate_FormURLEncoded(v *FormURLEncoded) error {
	errSet := github_com_go_courier_schema_pkg_validator.NewErrorSet()
	errSet.AddErr(validate_FormURLEncoded_Field_FormData(&v.FormData), github_com_go_courier_schema_pkg_validator.Location("body"), "")
	if errSet.Len() > 0 {
		return errSet
	}
	return nil

}
func validate_FormURLEncoded_Field_FormData(v *struct {
	String string   `name:"string"`
	Slice  []string `name:"slice"`
	Data   Data     `name:"data"`
}) error {
	errSet := github_com_go_courier_schema_pkg_validator.NewErrorSet()
	errSet.AddErr(validate_FormURLEncoded_Field_FormData_Field_Slice(&v.Slice), "slice")
	errSet.AddErr(validate_FormURLEncoded_Field_FormData_Field_Data(&v.Data), "data")
	if errSet.Len() > 0 {
		return errSet
	}
	return nil

}
func validate_FormURLEncoded_Field_FormData_Field_Slice(v *[]string) error {
	list := *v
	n := len(list)
	errSet := github_com_go_courier_schema_pkg_validator.NewErrorSet()
	for i := 0; i < n; i++ {
		errSet.AddErr(validate_FormURLEncoded_Field_FormData_Field_Slice_Elem(&list[i]), i)
	}
	if errSet.Len() > 0 {
		return errSet
	}
	return nil

}
func validate_FormURLEncoded_Field_FormData_Field_Slice_Elem(v *string) error {
	return nil
}
func validate_FormURLEncoded_Field_FormData_Field_Data(v *Data) error {
	return validate_Data(v)
}
func (v *FormURLEncoded) Validate() error {
	return validate_FormURLEncoded(v)
}
func validate_GetByID(v *GetByID) error {
	errSet := github_com_go_courier_schema_pkg_validator.NewErrorSet()
	errSet.AddErr(validate_GetByID_Field_Label(&v.Label), github_com_go_courier_schema_pkg_validator.Location("query"), "label")
	if errSet.Len() > 0 {
		return errSet
	}
	return nil

}
func validate_GetByID_Field_Label(v *[]string) error {
	list := *v
	n := len(list)
	errSet := github_com_go_courier_schema_pkg_validator.NewErrorSet()
	for i := 0; i < n; i++ {
		errSet.AddErr(validate_GetByID_Field_Label_Elem(&list[i]), i)
	}
	if errSet.Len() > 0 {
		return errSet
	}
	return nil

}
func validate_GetByID_Field_Label_Elem(v *string) error {
	return nil
}
func (v *GetByID) Validate() error {
	return validate_GetByID(v)
}
func validate_HealthCheck(v *HealthCheck) error {
	errSet := github_com_go_courier_schema_pkg_validator.NewErrorSet()
	if errSet.Len() > 0 {
		return errSet
	}
	return nil

}
func (v *HealthCheck) Validate() error {
	return validate_HealthCheck(v)
}
func validate_MustValidAccount(v *MustValidAccount) error {
	errSet := github_com_go_courier_schema_pkg_validator.NewErrorSet()
	if errSet.Len() > 0 {
		return errSet
	}
	return nil

}
func (v *MustValidAccount) Validate() error {
	return validate_MustValidAccount(v)
}
func validate_Proxy(v *Proxy) error {
	errSet := github_com_go_courier_schema_pkg_validator.NewErrorSet()
	if errSet.Len() > 0 {
		return errSet
	}
	return nil

}
func (v *Proxy) Validate() error {
	return validate_Proxy(v)
}
func validate_ProxyV2(v *ProxyV2) error {
	errSet := github_com_go_courier_schema_pkg_validator.NewErrorSet()
	if errSet.Len() > 0 {
		return errSet
	}
	return nil

}
func (v *ProxyV2) Validate() error {
	return validate_ProxyV2(v)
}
func validate_Redirect(v *Redirect) error {
	errSet := github_com_go_courier_schema_pkg_validator.NewErrorSet()
	if errSet.Len() > 0 {
		return errSet
	}
	return nil

}
func (v *Redirect) Validate() error {
	return validate_Redirect(v)
}
func validate_RedirectWhenError(v *RedirectWhenError) error {
	errSet := github_com_go_courier_schema_pkg_validator.NewErrorSet()
	if errSet.Len() > 0 {
		return errSet
	}
	return nil

}
func (v *RedirectWhenError) Validate() error {
	return validate_RedirectWhenError(v)
}
func validate_RemoveByID(v *RemoveByID) error {
	errSet := github_com_go_courier_schema_pkg_validator.NewErrorSet()
	if errSet.Len() > 0 {
		return errSet
	}
	return nil

}
func (v *RemoveByID) Validate() error {
	return validate_RemoveByID(v)
}
func validate_ShowImage(v *ShowImage) error {
	errSet := github_com_go_courier_schema_pkg_validator.NewErrorSet()
	if errSet.Len() > 0 {
		return errSet
	}
	return nil

}
func (v *ShowImage) Validate() error {
	return validate_ShowImage(v)
}
func validate_UpdateByID(v *UpdateByID) error {
	errSet := github_com_go_courier_schema_pkg_validator.NewErrorSet()
	errSet.AddErr(validate_UpdateByID_Field_Data(&v.Data), github_com_go_courier_schema_pkg_validator.Location("body"), "")
	if errSet.Len() > 0 {
		return errSet
	}
	return nil

}
func validate_UpdateByID_Field_Data(v *Data) error {
	return validate_Data(v)
}
func (v *UpdateByID) Validate() error {
	return validate_UpdateByID(v)
}
