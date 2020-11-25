package repository

func (r *Repository) Remove(name string) error {
	for i := range r.Dependencies {
		if r.Dependencies[i].Name == name {
			if err := r.Dependencies[i].Unlink(); err != nil {
				return err
			}

			r.Dependencies = append(r.Dependencies[:i], r.Dependencies[i+1:]...)

			return nil
		}
	}

	return ErrDependencyNotFound
}
