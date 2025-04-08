import { QueryClient, useMutation, useQuery } from "@tanstack/react-query"
import { createProduct, getProduct, getProducts } from "../services/products"
import { product } from "../apis/product"

const queryClient = new QueryClient()

export const useProduct = (id: number)=>{
    return useQuery({
        queryKey:["product", id],
        queryFn: ()=> getProduct(id),
    })
}

export const useProducts = () => {
    return useQuery({
        queryKey:["products"],
        queryFn: ()=> getProducts()
    })
}

export const useCreateProduct = (product: product.CreateProductRequest) => {
    return useMutation({
        mutationFn: () => createProduct(product),
        mutationKey: ["createProduct"],
        onSuccess: (data) => {
            // biome-ignore lint/suspicious/noConsole: <explanation>
            console.log("Product created successfully", data)
            queryClient.invalidateQueries({queryKey: ["products"]})
        },
        onError: (error) => {
            console.error("Error creating product", error)
        }
    })
}