import React, { useState } from 'react';

const server = process.env.REACT_APP_API_URL || 'http://127.0.0.1:9000';

interface Prop {
  onListingCompleted?: () => void;
}

interface ModalProps {
  onClose: () => void;
  children: React.ReactNode;
}

type formDataType = {
  name: string,
  category: string,
  image: string | File,
}

// Adding a modal window
const Modal: React.FC<ModalProps> = ({ onClose, children }) => {
  return (
    <div className="modalBackdrop" onClick={onClose}>
      <div className="modalContent" onClick={e => e.stopPropagation()}>
        {children}
        <button onClick={onClose} className="modalCloseButton">×</button>
      </div>
    </div>
  );
};

export const Listing: React.FC<Prop> = (props) => {
  const { onListingCompleted } = props;
  const initialState = {
    name: "",
    category: "",
    image: "",
  };
  const [values, setValues] = useState<formDataType>(initialState);
  const [categoryName, setCategoryName] = useState('');

  const onValueChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setValues({
      ...values, [event.target.name]: event.target.value,
    })
  };
  const onFileChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setValues({
      ...values, [event.target.name]: event.target.files![0],
    })
  };
  const onCategoryChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setCategoryName(event.target.value);
  };
  const onSubmit = (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault()
    const data = new FormData()
    data.append('name', values.name)
    data.append('category', values.category)
    data.append('image', values.image)

    fetch(server.concat('/items'), {
      method: 'POST',
      mode: 'cors',
      body: data,
    })
      .then(response => {
        console.log('POST status:', response.statusText);
        toggleItemFormVisibility();
        onListingCompleted && onListingCompleted();
      })
      .catch((error) => {
        console.error('POST error:', error);
      })
  };
  const onCategorySubmit = (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault()
    const data = new FormData()
    data.append('name', categoryName)

    fetch(server.concat('/categories'), {
      method: 'POST',
      mode: 'cors',
      body: data,
    })
      .then(response => {
        console.log('POST status:', response.statusText);
        toggleCategoryFormVisibility();
        onListingCompleted && onListingCompleted();
      })
      .catch((error) => {
        console.error('POST error:', error);
      })
  };
  const [isItemFormVisible, setIsItemFormVisible] = useState(false);
  const [isCategoryFormVisible, setIsCategoryFormVisible] = useState(false);
  const toggleItemFormVisibility = () => setIsItemFormVisible(!isItemFormVisible);
  const toggleCategoryFormVisibility = () => setIsCategoryFormVisible(!isCategoryFormVisible);
  return (
    <div className='Listing'>
      <div className="NavButton">
        <button onClick={toggleItemFormVisibility} className='CreateButton'>Create a item</button>
        <button onClick={toggleCategoryFormVisibility} className='CreateButton'>Create a Category</button>
      </div>
      {isItemFormVisible && (
        <Modal onClose={toggleItemFormVisibility}>
          <form onSubmit={onSubmit}>
            <label htmlFor="name" className='ItemName'>Name:</label>
            <input type='text' name='name' id='name' placeholder='name' onChange={onValueChange} required />
            <label htmlFor="name" className='ItemCategory'>Category:</label>
            <input type='text' name='category' id='category' placeholder='category' onChange={onValueChange} required />
            <label htmlFor="name" className='ItemImage'>Image:</label>
            <input type='file' name='image' id='image' onChange={onFileChange} required />
            <button type='submit'>List this item</button>
          </form>
        </Modal>
      )}

      {isCategoryFormVisible && (
        <Modal onClose={toggleCategoryFormVisibility}>
          <form onSubmit={onCategorySubmit}>
          <label htmlFor="categoryName" className='CategoryName'>Category Name:</label>
            <input type='text' name='categoryName' id='categoryName' placeholder='new category' onChange={onCategoryChange} required />
            <button type='submit'>Add a Category</button>
          </form>
        </Modal>
      )}
    </div>
  );
}
