import { product } from "../apis/product"
import api from "../utils/api"

export const getProduct = (id: number) =>{
    return api.get<product.ProductResponse>(`/products/${id}`).then((res) => {
        return res.data
    })
}

export const getProducts = () => {
    return api.get<product.ListProductResponse>('/products').then((res) => {
        return res.data
    })
}

export const createProduct = (product: product.CreateProductRequest)=>{
    return api.post<product.ProductResponse>('/products', product).then((res) => {
        return res.data
    })
}