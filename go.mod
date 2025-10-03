module github.com/agumiroff/BigTechProject

go 1.24.4

replace (
	github.com/agumiroff/BigTechProject/inventory => ./inventory
	github.com/agumiroff/BigTechProject/order => ./order
	github.com/agumiroff/BigTechProject/payment => ./payment
	github.com/agumiroff/BigTechProject/shared => ./shared
)
