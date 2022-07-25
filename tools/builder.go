package tools

func NewUgo(name string) *UGO {
	return &UGO{
		Layout: 1,
		Title: UgoTitle{
			Names: [5]string{name, name, name, name, name},
		},
		Assets: []UgoAsset{},
	}
}

func (u *UGO) AddButton(name string, imageIndex int) error {
	buttonAsset := UgoAsset{
		Type:   TypeButton,
		Name:   name,
		Index1: imageIndex,
	}

	u.Assets = append(u.Assets, buttonAsset)
	return nil
}

func (u *UGO) AddCategory(name string, selected bool) error {
	selectedNum := 0
	if selected {
		selectedNum = 1
	}

	categoryAsset := UgoAsset{
		Type:   TypeCategory,
		Name:   name,
		Index1: selectedNum,
	}

	u.Assets = append(u.Assets, categoryAsset)
	return nil
}
