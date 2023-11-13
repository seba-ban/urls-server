function cellEditedFactory(propertyName) {
  return async function (cell) {
    const url = cell.getData().url;
    const prop = cell.getData()[propertyName];
    await fetch("/urls", {
      method: "PATCH",
      body: JSON.stringify({
        url,
        [propertyName]: prop,
      }),
      headers: {
        "Content-type": "application/json; charset=UTF-8",
      },
    });
  };
}

const commonOpts = {
  paginationSize: 15,
  pagination: true,
  // maxHeight: "50vh",
  layout: "fitColumns",
  ajaxURL: "/urls",
};

const columns = [
  {
    title: "URL",
    field: "url",
    width: 350,
    headerFilter: true,
    formatter: "link",
    formatterParams: {
      target: "_blank",
    },
  },
  {
    title: "Description",
    field: "description",
    headerFilter: true,
  },
  {
    title: "Created At",
    field: "created_at",
    hozAlign: "center",
    width: 150,
  },
];

const doneTable = new Tabulator("#done", {
  // data: done, //assign data to table
  ...commonOpts,
  ajaxParams: { pending: false },
  columns: [
    ...columns,
    {
      title: "Done At",
      field: "read_at",
      hozAlign: "center",
      width: 150,
    },
  ],
});

const pendingTable = new Tabulator("#table", {
  // data: pending, //assign data to table
  ...commonOpts,
  ajaxParams: { pending: true },
  columns: [
    ...columns,
    {
      title: "Priority",
      field: "priority",
      hozAlign: "center",
      width: 50,
      cellEdited: cellEditedFactory("priority"),
      sorter: "number",
      editor: "number",
      editorParams: {
        step: 1,
      },
    },
    {
      title: "Done",
      field: "read_at",
      formatter: "buttonCross",
      width: 50,
      hozAlign: "center",
      cellEdited: (cell) => {
        cellEditedFactory("read_at")(cell);
        const data = cell.getData();
        // TODO: bleh
        data.read_at = new Date().toISOString();
        cell.getRow().delete();
        doneTable.addRow(data);
      },
      editor: "tickCross",
    },
  ],
  initialSort: [
    { column: "created_at", dir: "desc" },
    { column: "priority", dir: "desc" },
  ],
});
