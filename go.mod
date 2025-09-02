module github.com/agumiroff/BigTechProject

go 1.24.4

require (
    github.com/agumiroff/BigTechProject/shared v0.0.0
    github.com/agumiroff/BigTechProject/payment v0.0.0
    github.com/agumiroff/BigTechProject/order v0.0.0
    github.com/agumiroff/BigTechProject/inventory v0.0.0
)

replace (
    github.com/agumiroff/BigTechProject/shared => ./shared
    github.com/agumiroff/BigTechProject/payment => ./payment
    github.com/agumiroff/BigTechProject/order => ./order
    github.com/agumiroff/BigTechProject/inventory => ./inventory
)
