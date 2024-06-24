package main

type Updater interface {
	Update()
}

type UpdateHandler struct {
	updaters []Updater
}

func (u *UpdateHandler) Add(updater Updater) {
	u.updaters = append(u.updaters, updater)
}

func (u *UpdateHandler) HandleUpdate() {
	for _, updater := range u.updaters {
		updater.Update()
	}
}

func (u *UpdateHandler) Remove(updater Updater) {
	for i, v := range u.updaters {
		if v == updater {
			u.updaters = append(u.updaters[:i], u.updaters[i+1:]...)
			return
		}
	}
}

func (u *UpdateHandler) Clear() {
	u.updaters = nil
}
