import React, {useState}  from 'react'
import {Link} from "react-router-dom";
import ErrorBox from './ErrorBox'
import axios from "axios";

function Create() {

  const [title, setTitle] = useState("");
  const [shortDescription, setShortDescription] = useState("");
  const [fullDescription, setFullDescription] = useState("");
  const [imageFile, setImageFile] = useState("");
  const [fileType, setFileType] = useState("");
  const [price, setPrice] = useState(0);
  const [dataset, setDataset] = useState("");
  const [isError, setIsError] = useState(false);
  const [error, setError] = useState(false);

  const convertBase64 = (file) => {
    return new Promise((resolve, reject) => {
      const fileReader = new FileReader();
      fileReader.readAsBinaryString(file)
      fileReader.onload = () => {
        resolve(fileReader.result);
      }
      fileReader.onerror = (error) => {
        reject(error);
      }
    })
  }

  const HandleFormSubmit = async (e) => {
    e.preventDefault();

    // Convert image file to base64 string
    const blobString = await convertBase64(imageFile);
    const fileString = btoa(blobString);

    const data = {
      title: title,
      shortDescription: shortDescription,
      fullDescription: fullDescription,
      image: fileString,
      fileType: fileType,
      price: Number(price),
    };

    console.log(data);

    const formData = new FormData();
    formData.append('metadata', JSON.stringify(data));
    formData.append('file', fileString);

    const csrftoken = localStorage.getItem('token');
    const instance = axios.create({
      baseURL: "",
      headers: {
        "x-csrf-token": csrftoken,
        "content-type": "multipart/form-data"
      }
    })

    const url = "/api/v1/dataset";
    const apiReq = await instance.post(
        url,
        formData
    );
    console.log(apiReq);

    return false;
  };

  const HandleDatasetImage = (e) => {
    setImageFile(e.target.files[0]);
  }

  const HandleDataset = (e) => {
    setDataset(e.target.files[0]);
  }

  return (
    <div class="CreateDataset">
      <h2>Create dataset</h2>
      <form onSubmit={HandleFormSubmit}>
        <label>
          Title*
          <div>
            <input type="text" name="title" placeholder="Title" onChange={e => setTitle(e.target.value)}/>
            <span>Set a clear description title for your dataset.</span>
          </div>
        </label>

        <label>
          Short description*
          <div>
            <input type="text" name="shortDescription" placeholder="(100 char max)" onChange={e => setShortDescription(e.target.value)}/>
            <span>Explain your dataset in 50 characters or less.</span>
          </div>
        </label>

        <label>
          Full description*
          <div>
            <textarea name="fullDescription" placeholder="Enter description" onChange={e => setFullDescription(e.target.value)}/>
            <span>Fully describe the contents in the dataset and provide example of how the data is structured. The more information the better.</span>
          </div>
        </label>

        <label>
          Image*
          <div>
            <input type="file" name="imageFile" onChange={HandleDatasetImage}/>
            <span>Attach a JPG or PNG cover photo for your dataset.</span>
          </div>
        </label>

        <div class="form-divider"></div>

        <h3>Files</h3>

        <label>
          File Type*
          <div>
            <input type="text" name="filetype" placeholder="CSV, Excel, SQL, MP4, etc" onChange={e => setFileType(e.target.value)}/>
            <span>Specify the file type for the dataset.</span>
          </div>
        </label>

        <label>
          Price*
          <div>
            <input type="text" name="price" placeholder="5.23" onChange={e => setPrice(e.target.value)}/>
            <span>Set your price in Filecoin (FIL).</span>
          </div>
        </label>

        <label>
          Dataset*
          <div>
            <input type="file" name="dataset" onChange={HandleDataset}/>
            <span>Finally, attach your dataset file(s) in zip format.</span>
          </div>
        </label>

        <div className="form-divider"></div>

        <div>
          <input type="submit" value="Submit" className="orange-button"/>
        </div>

      </form>
    </div>
  )
}

export default Create