import React, { useState } from 'react';
import './App.css';
import { useCreateProduct, useProducts } from './queries/product';
import { product } from './apis/product';

function App() {
  const { data, isLoading } = useProducts();
  const [state, setState] = useState({
    name: '',
    description: '',
    price: 0,
  });

  const createProductMutate = useCreateProduct(
    state as product.CreateProductRequest
  );

  const handleCreateProduct = () => {
    createProductMutate.mutate();
  };

  const handleChange = (
    e: React.ChangeEvent<HTMLInputElement>,
    type: string
  ) => {
    setState((prev) => ({
      ...prev,
      [type]: type === 'price' ? Number(e.target.value) : e.target.value,
    }));
  };

  return (
    <>
      <div>
        <input
          placeholder="name"
          value={state.name}
          onChange={(e) => handleChange(e, 'name')}
        />
        <input
          placeholder="description"
          value={state.description}
          onChange={(e) => handleChange(e, 'description')}
        />
        <input
          placeholder="price"
          type="number"
          value={state.price}
          onChange={(e) => handleChange(e, 'price')}
        />
        <button type="submit" onClick={handleCreateProduct}>
          Create
        </button>
      </div>
      <div>
        {isLoading ? (
          <div>Loading...</div>
        ) : (
          data?.products?.map((product) => (
            <div key={product.id}>
              <h2>{product.name}</h2>
              <p>{product.description}</p>
              <p>{product.price}</p>
            </div>
          ))
        )}
      </div>
    </>
  );
}

export default App;
