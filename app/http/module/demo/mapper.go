package demo

import demoService "github.com/YunzeGao/fire/app/provider/demo"

func UserModelsToUserDTOs(models []UserModel) []UserDTO {
	var ret []UserDTO
	for _, model := range models {
		t := UserDTO{
			ID:   model.UserId,
			Name: model.Name,
		}
		ret = append(ret, t)
	}
	return ret
}

func StudentsToUserDTOs(students []demoService.Student) []UserDTO {
	var ret []UserDTO
	for _, student := range students {
		t := UserDTO{
			ID:   student.ID,
			Name: student.Name,
		}
		ret = append(ret, t)
	}
	return ret
}
