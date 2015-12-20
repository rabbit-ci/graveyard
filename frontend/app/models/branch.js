import DS from 'ember-data';

export default DS.Model.extend({
  name: DS.attr('string'),
  existsInGit: DS.attr('boolean'),
  project: DS.belongsTo('project'),
  builds: DS.hasMany('build', { async: true })
});
