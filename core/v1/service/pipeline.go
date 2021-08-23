package service

type Pipeline interface {
	Apply(url,revision string)error
}