import React from 'react'
import {Link} from "react-router-dom";
import ErrorBox from './ErrorBox'

function Create() {
  return (
    <div class="CreateDataset">
      <h2>Create dataset</h2>
      <form>
        <label>
          Title*
          <div>
            <input type="text" name="title" placeholder="Title" />
            <span>Set a clear description title for your dataset.</span>
          </div>
        </label>

        <label>
          Short description*
          <div>
            <input type="text" name="shortdescription" placeholder="(100 char max))" />
            <span>Explain your dataset in 50 characters or less.</span>
          </div>
        </label>

        <label>
          Full description*
          <div>
            <textarea name="fulldescription" placeholder="Enter description" />
            <span>Fully describe the contents in the dataset and provide example of how the data is structured. The more information the better.</span>
          </div>
        </label>

        <label>
          Image*
          <div>
            <input type="file"/>
            <span>Attach a JPG or PNG cover photo for your dataset.</span>
          </div>
        </label>

        <div class="form-divider"></div>

        <h3>Files</h3>

        <label>
          File Type*
          <div>
            <input type="text" name="filetype" placeholder="CSV, Excel, SQL, MP4, etc" />
            <span>Specify the file type for the dataset.</span>
          </div>
        </label>

        <label>
          Price*
          <div>
            <input type="text" name="price" placeholder="5.23" />
            <span>Set your price in Filecoin (FIL).</span>
          </div>
        </label>

        <label>
          Dataset*
          <div>
            <input type="file" name="dataset" placeholder="Title" />
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